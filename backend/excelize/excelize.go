// Copyright 2016 - 2020 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package excelize providing a set of functions that allow you to write to
// and read from XLSX / XLSM / XLTM files. Supports reading and writing
// spreadsheet documents generated by Microsoft Exce™ 2007 and later. Supports
// complex components by high compatibility, and provided streaming API for
// generating or reading data from a worksheet with huge amounts of data. This
// library needs Go version 1.10 or later.
//
// See https://xuri.me/excelize for more information about this package.
package excelize

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html/charset"
)

// File define a populated spreadsheet file struct.
type File struct {
	sync.Mutex
	options          *Options
	xmlAttr          map[string][]xml.Attr
	checked          map[string]bool
	sheetMap         map[string]string
	streams          map[string]*StreamWriter
	CalcChain        *xlsxCalcChain
	Comments         map[string]*xlsxComments
	ContentTypes     *xlsxTypes
	Drawings         map[string]*xlsxWsDr
	Path             string
	SharedStrings    *xlsxSST
	sharedStringsMap map[string]int
	Sheet            map[string]*xlsxWorksheet
	SheetCount       int
	Styles           *xlsxStyleSheet
	Theme            *xlsxTheme
	DecodeVMLDrawing map[string]*decodeVmlDrawing
	VMLDrawing       map[string]*vmlDrawing
	WorkBook         *xlsxWorkbook
	Relationships    map[string]*xlsxRelationships
	XLSX             map[string][]byte
	CharsetReader    charsetTranscoderFn
}

type charsetTranscoderFn func(charset string, input io.Reader) (rdr io.Reader, err error)

// Options define the options for open spreadsheet.
type Options struct {
	Password string
}

// OpenFile take the name of an spreadsheet file and returns a populated spreadsheet file struct
// for it. For example, open spreadsheet with password protection:
//
//    f, err := excelize.OpenFile("Book1.xlsx", excelize.Options{Password: "password"})
//    if err != nil {
//        return
//    }
//
// Note that the excelize just support decrypt and not support encrypt currently, the spreadsheet
// saved by Save and SaveAs will be without password unprotected.
func OpenFile(filename string, opt ...Options) (*File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	f, err := OpenReader(file, opt...)
	if err != nil {
		return nil, err
	}
	f.Path = filename
	return f, nil
}

// newFile is object builder
func newFile() *File {
	return &File{
		xmlAttr:          make(map[string][]xml.Attr),
		checked:          make(map[string]bool),
		sheetMap:         make(map[string]string),
		Comments:         make(map[string]*xlsxComments),
		Drawings:         make(map[string]*xlsxWsDr),
		sharedStringsMap: make(map[string]int),
		Sheet:            make(map[string]*xlsxWorksheet),
		DecodeVMLDrawing: make(map[string]*decodeVmlDrawing),
		VMLDrawing:       make(map[string]*vmlDrawing),
		Relationships:    make(map[string]*xlsxRelationships),
		CharsetReader:    charset.NewReaderLabel,
	}
}

// OpenReader read data stream from io.Reader and return a populated
// spreadsheet file.
func OpenReader(r io.Reader, opt ...Options) (*File, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f := newFile()
	if bytes.Contains(b, oleIdentifier) && len(opt) > 0 {
		for _, o := range opt {
			f.options = &o
		}
		b, err = Decrypt(b, f.options)
		if err != nil {
			return nil, fmt.Errorf("decrypted file failed")
		}
	}
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}

	file, sheetCount, err := ReadZipReader(zr)
	if err != nil {
		return nil, err
	}
	f.SheetCount, f.XLSX = sheetCount, file
	f.CalcChain = f.calcChainReader()
	f.sheetMap = f.getSheetMap()
	f.Styles = f.stylesReader()
	f.Theme = f.themeReader()
	return f, nil
}

// CharsetTranscoder Set user defined codepage transcoder function for open
// XLSX from non UTF-8 encoding.
func (f *File) CharsetTranscoder(fn charsetTranscoderFn) *File { f.CharsetReader = fn; return f }

// Creates new XML decoder with charset reader.
func (f *File) xmlNewDecoder(rdr io.Reader) (ret *xml.Decoder) {
	ret = xml.NewDecoder(rdr)
	ret.CharsetReader = f.CharsetReader
	return
}

// setDefaultTimeStyle provides a function to set default numbers format for
// time.Time type cell value by given worksheet name, cell coordinates and
// number format code.
func (f *File) setDefaultTimeStyle(sheet, axis string, format int) error {
	s, err := f.GetCellStyle(sheet, axis)
	if err != nil {
		return err
	}
	if s == 0 {
		style, _ := f.NewStyle(&Style{NumFmt: format})
		err = f.SetCellStyle(sheet, axis, axis, style)
	}
	return err
}

// workSheetReader provides a function to get the pointer to the structure
// after deserialization by given worksheet name.
func (f *File) workSheetReader(sheet string) (ws *xlsxWorksheet, err error) {
	f.Lock()
	defer f.Unlock()
	var (
		name string
		ok   bool
	)

	if name, ok = f.sheetMap[trimSheetName(sheet)]; !ok {
		err = fmt.Errorf("sheet %s is not exist", sheet)
		return
	}
	if ws = f.Sheet[name]; f.Sheet[name] == nil {
		if strings.HasPrefix(name, "xl/chartsheets") {
			err = fmt.Errorf("sheet %s is chart sheet", sheet)
			return
		}
		ws = new(xlsxWorksheet)
		if _, ok := f.xmlAttr[name]; !ok {
			d := f.xmlNewDecoder(bytes.NewReader(namespaceStrictToTransitional(f.readXML(name))))
			f.xmlAttr[name] = append(f.xmlAttr[name], getRootElement(d)...)
		}
		if err = f.xmlNewDecoder(bytes.NewReader(namespaceStrictToTransitional(f.readXML(name)))).
			Decode(ws); err != nil && err != io.EOF {
			err = fmt.Errorf("xml decode error: %s", err)
			return
		}
		err = nil
		if f.checked == nil {
			f.checked = make(map[string]bool)
		}
		if ok = f.checked[name]; !ok {
			checkSheet(ws)
			if err = checkRow(ws); err != nil {
				return
			}
			f.checked[name] = true
		}
		f.Sheet[name] = ws
	}

	return
}

// checkSheet provides a function to fill each row element and make that is
// continuous in a worksheet of XML.
func checkSheet(ws *xlsxWorksheet) {
	var row int
	for _, r := range ws.SheetData.Row {
		if r.R != 0 && r.R > row {
			row = r.R
			continue
		}
		if r.R != row {
			row++
		}
	}
	sheetData := xlsxSheetData{Row: make([]xlsxRow, row)}
	row = 0
	for _, r := range ws.SheetData.Row {
		if r.R == row {
			sheetData.Row[r.R-1].C = append(sheetData.Row[r.R-1].C, r.C...)
			continue
		}
		if r.R != 0 {
			sheetData.Row[r.R-1] = r
			row = r.R
			continue
		}
		row++
		r.R = row
		sheetData.Row[row-1] = r
	}
	for i := 1; i <= row; i++ {
		sheetData.Row[i-1].R = i
	}
	ws.SheetData = sheetData
}

// addRels provides a function to add relationships by given XML path,
// relationship type, target and target mode.
func (f *File) addRels(relPath, relType, target, targetMode string) int {
	var uniqPart = map[string]string{
		SourceRelationshipSharedStrings: "/xl/sharedStrings.xml",
	}
	rels := f.relsReader(relPath)
	if rels == nil {
		rels = &xlsxRelationships{}
	}
	var rID int
	for idx, rel := range rels.Relationships {
		ID, _ := strconv.Atoi(strings.TrimPrefix(rel.ID, "rId"))
		if ID > rID {
			rID = ID
		}
		if relType == rel.Type {
			if partName, ok := uniqPart[rel.Type]; ok {
				rels.Relationships[idx].Target = partName
				return rID
			}
		}
	}
	rID++
	var ID bytes.Buffer
	ID.WriteString("rId")
	ID.WriteString(strconv.Itoa(rID))
	rels.Relationships = append(rels.Relationships, xlsxRelationship{
		ID:         ID.String(),
		Type:       relType,
		Target:     target,
		TargetMode: targetMode,
	})
	f.Relationships[relPath] = rels
	return rID
}

// UpdateLinkedValue fix linked values within a spreadsheet are not updating in
// Office Excel 2007 and 2010. This function will be remove value tag when met a
// cell have a linked value. Reference
// https://social.technet.microsoft.com/Forums/office/en-US/e16bae1f-6a2c-4325-8013-e989a3479066/excel-2010-linked-cells-not-updating
//
// Notice: after open XLSX file Excel will be update linked value and generate
// new value and will prompt save file or not.
//
// For example:
//
//    <row r="19" spans="2:2">
//        <c r="B19">
//            <f>SUM(Sheet2!D2,Sheet2!D11)</f>
//            <v>100</v>
//         </c>
//    </row>
//
// to
//
//    <row r="19" spans="2:2">
//        <c r="B19">
//            <f>SUM(Sheet2!D2,Sheet2!D11)</f>
//        </c>
//    </row>
//
func (f *File) UpdateLinkedValue() error {
	wb := f.workbookReader()
	// recalculate formulas
	wb.CalcPr = nil
	for _, name := range f.GetSheetList() {
		xlsx, err := f.workSheetReader(name)
		if err != nil {
			return err
		}
		for indexR := range xlsx.SheetData.Row {
			for indexC, col := range xlsx.SheetData.Row[indexR].C {
				if col.F != nil && col.V != "" {
					xlsx.SheetData.Row[indexR].C[indexC].V = ""
					xlsx.SheetData.Row[indexR].C[indexC].T = ""
				}
			}
		}
	}
	return nil
}

// AddVBAProject provides the method to add vbaProject.bin file which contains
// functions and/or macros. The file extension should be .xlsm. For example:
//
//    if err := f.SetSheetPrOptions("Sheet1", excelize.CodeName("Sheet1")); err != nil {
//        fmt.Println(err)
//    }
//    if err := f.AddVBAProject("vbaProject.bin"); err != nil {
//        fmt.Println(err)
//    }
//    if err := f.SaveAs("macros.xlsm"); err != nil {
//        fmt.Println(err)
//    }
//
func (f *File) AddVBAProject(bin string) error {
	var err error
	// Check vbaProject.bin exists first.
	if _, err = os.Stat(bin); os.IsNotExist(err) {
		return fmt.Errorf("stat %s: no such file or directory", bin)
	}
	if path.Ext(bin) != ".bin" {
		return errors.New("unsupported VBA project extension")
	}
	f.setContentTypePartVBAProjectExtensions()
	wb := f.relsReader(f.getWorkbookRelsPath())
	var rID int
	var ok bool
	for _, rel := range wb.Relationships {
		if rel.Target == "vbaProject.bin" && rel.Type == SourceRelationshipVBAProject {
			ok = true
			continue
		}
		t, _ := strconv.Atoi(strings.TrimPrefix(rel.ID, "rId"))
		if t > rID {
			rID = t
		}
	}
	rID++
	if !ok {
		wb.Relationships = append(wb.Relationships, xlsxRelationship{
			ID:     "rId" + strconv.Itoa(rID),
			Target: "vbaProject.bin",
			Type:   SourceRelationshipVBAProject,
		})
	}
	file, _ := ioutil.ReadFile(bin)
	f.XLSX["xl/vbaProject.bin"] = file
	return err
}

// setContentTypePartVBAProjectExtensions provides a function to set the
// content type for relationship parts and the main document part.
func (f *File) setContentTypePartVBAProjectExtensions() {
	var ok bool
	content := f.contentTypesReader()
	for _, v := range content.Defaults {
		if v.Extension == "bin" {
			ok = true
		}
	}
	for idx, o := range content.Overrides {
		if o.PartName == "/xl/workbook.xml" {
			content.Overrides[idx].ContentType = ContentTypeMacro
		}
	}
	if !ok {
		content.Defaults = append(content.Defaults, xlsxDefault{
			Extension:   "bin",
			ContentType: ContentTypeVBA,
		})
	}
}
