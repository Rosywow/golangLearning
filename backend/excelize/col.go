// Copyright 2016 - 2020 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to
// and read from XLSX / XLSM / XLTM files. Supports reading and writing
// spreadsheet documents generated by Microsoft Exce™ 2007 and later. Supports
// complex components by high compatibility, and provided streaming API for
// generating or reading data from a worksheet with huge amounts of data. This
// library needs Go version 1.10 or later.

package excelize

import (
	"bytes"
	"encoding/xml"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/mohae/deepcopy"
)

// Define the default cell size and EMU unit of measurement.
const (
	defaultColWidthPixels  float64 = 64
	defaultRowHeight       float64 = 15
	defaultRowHeightPixels float64 = 20
	EMU                    int     = 9525
)

// Cols defines an iterator to a sheet
type Cols struct {
	err                                  error
	curCol, totalCol, stashCol, totalRow int
	sheet                                string
	cols                                 []xlsxCols
	f                                    *File
	sheetXML                             []byte
}

// GetCols return all the columns in a sheet by given worksheet name (case
// sensitive). For example:
//
//    cols, err := f.GetCols("Sheet1")
//    if err != nil {
//        fmt.Println(err)
//        return
//    }
//    for _, col := range cols {
//        for _, rowCell := range col {
//            fmt.Print(rowCell, "\t")
//        }
//        fmt.Println()
//    }
//
func (f *File) GetCols(sheet string) ([][]string, error) {
	cols, err := f.Cols(sheet)
	if err != nil {
		return nil, err
	}
	results := make([][]string, 0, 64)
	for cols.Next() {
		col, _ := cols.Rows()
		results = append(results, col)
	}
	return results, nil
}

// Next will return true if the next column is found.
func (cols *Cols) Next() bool {
	cols.curCol++
	return cols.curCol <= cols.totalCol
}

// Error will return an error when the error occurs.
func (cols *Cols) Error() error {
	return cols.err
}

// Rows return the current column's row values.
func (cols *Cols) Rows() ([]string, error) {
	var (
		err              error
		inElement        string
		cellCol, cellRow int
		rows             []string
	)
	if cols.stashCol >= cols.curCol {
		return rows, err
	}
	d := cols.f.sharedStringsReader()
	decoder := cols.f.xmlNewDecoder(bytes.NewReader(cols.sheetXML))
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			inElement = startElement.Name.Local
			if inElement == "row" {
				cellCol = 0
				cellRow++
				attrR, _ := attrValToInt("r", startElement.Attr)
				if attrR != 0 {
					cellRow = attrR
				}
			}
			if inElement == "c" {
				cellCol++
				for _, attr := range startElement.Attr {
					if attr.Name.Local == "r" {
						if cellCol, cellRow, err = CellNameToCoordinates(attr.Value); err != nil {
							return rows, err
						}
					}
				}
				blank := cellRow - len(rows)
				for i := 1; i < blank; i++ {
					rows = append(rows, "")
				}
				if cellCol == cols.curCol {
					colCell := xlsxC{}
					_ = decoder.DecodeElement(&colCell, &startElement)
					val, _ := colCell.getValueFrom(cols.f, d)
					rows = append(rows, val)
				}
			}
		}
	}
	return rows, nil
}

// Cols returns a columns iterator, used for streaming reading data for a
// worksheet with a large data. For example:
//
//    cols, err := f.Cols("Sheet1")
//    if err != nil {
//        fmt.Println(err)
//        return
//    }
//    for cols.Next() {
//        col, err := cols.Rows()
//        if err != nil {
//            fmt.Println(err)
//        }
//        for _, rowCell := range col {
//            fmt.Print(rowCell, "\t")
//        }
//        fmt.Println()
//    }
//
func (f *File) Cols(sheet string) (*Cols, error) {
	name, ok := f.sheetMap[trimSheetName(sheet)]
	if !ok {
		return nil, ErrSheetNotExist{sheet}
	}
	if f.Sheet[name] != nil {
		output, _ := xml.Marshal(f.Sheet[name])
		f.saveFileList(name, f.replaceNameSpaceBytes(name, output))
	}
	var (
		inElement            string
		cols                 Cols
		cellCol, curRow, row int
		err                  error
	)
	cols.sheetXML = f.readXML(name)
	decoder := f.xmlNewDecoder(bytes.NewReader(cols.sheetXML))
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			inElement = startElement.Name.Local
			if inElement == "row" {
				row++
				for _, attr := range startElement.Attr {
					if attr.Name.Local == "r" {
						if curRow, err = strconv.Atoi(attr.Value); err != nil {
							return &cols, err
						}
						row = curRow
					}
				}
				cols.totalRow = row
				cellCol = 0
			}
			if inElement == "c" {
				cellCol++
				for _, attr := range startElement.Attr {
					if attr.Name.Local == "r" {
						if cellCol, _, err = CellNameToCoordinates(attr.Value); err != nil {
							return &cols, err
						}
					}
				}
				if cellCol > cols.totalCol {
					cols.totalCol = cellCol
				}
			}
		}
	}
	cols.f = f
	cols.sheet = trimSheetName(sheet)
	return &cols, nil
}

// GetColVisible provides a function to get visible of a single column by given
// worksheet name and column name. For example, get visible state of column D
// in Sheet1:
//
//    visible, err := f.GetColVisible("Sheet1", "D")
//
func (f *File) GetColVisible(sheet, col string) (bool, error) {
	visible := true
	colNum, err := ColumnNameToNumber(col)
	if err != nil {
		return visible, err
	}

	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return false, err
	}
	if ws.Cols == nil {
		return visible, err
	}

	for c := range ws.Cols.Col {
		colData := &ws.Cols.Col[c]
		if colData.Min <= colNum && colNum <= colData.Max {
			visible = !colData.Hidden
		}
	}
	return visible, err
}

// SetColVisible provides a function to set visible columns by given worksheet
// name, columns range and visibility.
//
// For example hide column D on Sheet1:
//
//    err := f.SetColVisible("Sheet1", "D", false)
//
// Hide the columns from D to F (included):
//
//    err := f.SetColVisible("Sheet1", "D:F", false)
//
func (f *File) SetColVisible(sheet, columns string, visible bool) error {
	var max int

	colsTab := strings.Split(columns, ":")
	min, err := ColumnNameToNumber(colsTab[0])
	if err != nil {
		return err
	}
	if len(colsTab) == 2 {
		max, err = ColumnNameToNumber(colsTab[1])
		if err != nil {
			return err
		}
	} else {
		max = min
	}
	if max < min {
		min, max = max, min
	}
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	colData := xlsxCol{
		Min:         min,
		Max:         max,
		Width:       9, // default width
		Hidden:      !visible,
		CustomWidth: true,
	}
	if ws.Cols == nil {
		cols := xlsxCols{}
		cols.Col = append(cols.Col, colData)
		ws.Cols = &cols
		return nil
	}
	ws.Cols.Col = flatCols(colData, ws.Cols.Col, func(fc, c xlsxCol) xlsxCol {
		fc.BestFit = c.BestFit
		fc.Collapsed = c.Collapsed
		fc.CustomWidth = c.CustomWidth
		fc.OutlineLevel = c.OutlineLevel
		fc.Phonetic = c.Phonetic
		fc.Style = c.Style
		fc.Width = c.Width
		return fc
	})
	return nil
}

// GetColOutlineLevel provides a function to get outline level of a single
// column by given worksheet name and column name. For example, get outline
// level of column D in Sheet1:
//
//    level, err := f.GetColOutlineLevel("Sheet1", "D")
//
func (f *File) GetColOutlineLevel(sheet, col string) (uint8, error) {
	level := uint8(0)
	colNum, err := ColumnNameToNumber(col)
	if err != nil {
		return level, err
	}
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return 0, err
	}
	if ws.Cols == nil {
		return level, err
	}
	for c := range ws.Cols.Col {
		colData := &ws.Cols.Col[c]
		if colData.Min <= colNum && colNum <= colData.Max {
			level = colData.OutlineLevel
		}
	}
	return level, err
}

// SetColOutlineLevel provides a function to set outline level of a single
// column by given worksheet name and column name. The value of parameter
// 'level' is 1-7. For example, set outline level of column D in Sheet1 to 2:
//
//    err := f.SetColOutlineLevel("Sheet1", "D", 2)
//
func (f *File) SetColOutlineLevel(sheet, col string, level uint8) error {
	if level > 7 || level < 1 {
		return errors.New("invalid outline level")
	}
	colNum, err := ColumnNameToNumber(col)
	if err != nil {
		return err
	}
	colData := xlsxCol{
		Min:          colNum,
		Max:          colNum,
		OutlineLevel: level,
		CustomWidth:  true,
	}
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	if ws.Cols == nil {
		cols := xlsxCols{}
		cols.Col = append(cols.Col, colData)
		ws.Cols = &cols
		return err
	}
	ws.Cols.Col = flatCols(colData, ws.Cols.Col, func(fc, c xlsxCol) xlsxCol {
		fc.BestFit = c.BestFit
		fc.Collapsed = c.Collapsed
		fc.CustomWidth = c.CustomWidth
		fc.Hidden = c.Hidden
		fc.Phonetic = c.Phonetic
		fc.Style = c.Style
		fc.Width = c.Width
		return fc
	})
	return err
}

// SetColStyle provides a function to set style of columns by given worksheet
// name, columns range and style ID.
//
// For example set style of column H on Sheet1:
//
//    err = f.SetColStyle("Sheet1", "H", style)
//
// Set style of columns C:F on Sheet1:
//
//    err = f.SetColStyle("Sheet1", "C:F", style)
//
func (f *File) SetColStyle(sheet, columns string, styleID int) error {
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	var c1, c2 string
	var min, max int
	cols := strings.Split(columns, ":")
	c1 = cols[0]
	min, err = ColumnNameToNumber(c1)
	if err != nil {
		return err
	}
	if len(cols) == 2 {
		c2 = cols[1]
		max, err = ColumnNameToNumber(c2)
		if err != nil {
			return err
		}
	} else {
		max = min
	}
	if max < min {
		min, max = max, min
	}
	if ws.Cols == nil {
		ws.Cols = &xlsxCols{}
	}
	ws.Cols.Col = flatCols(xlsxCol{
		Min:   min,
		Max:   max,
		Width: 9,
		Style: styleID,
	}, ws.Cols.Col, func(fc, c xlsxCol) xlsxCol {
		fc.BestFit = c.BestFit
		fc.Collapsed = c.Collapsed
		fc.CustomWidth = c.CustomWidth
		fc.Hidden = c.Hidden
		fc.OutlineLevel = c.OutlineLevel
		fc.Phonetic = c.Phonetic
		fc.Width = c.Width
		return fc
	})
	return nil
}

// SetColWidth provides a function to set the width of a single column or
// multiple columns. For example:
//
//    f := excelize.NewFile()
//    err := f.SetColWidth("Sheet1", "A", "H", 20)
//
func (f *File) SetColWidth(sheet, startcol, endcol string, width float64) error {
	min, err := ColumnNameToNumber(startcol)
	if err != nil {
		return err
	}
	max, err := ColumnNameToNumber(endcol)
	if err != nil {
		return err
	}
	if width > MaxColumnWidth {
		return errors.New("the width of the column must be smaller than or equal to 255 characters")
	}
	if min > max {
		min, max = max, min
	}

	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	col := xlsxCol{
		Min:         min,
		Max:         max,
		Width:       width,
		CustomWidth: true,
	}
	if ws.Cols == nil {
		cols := xlsxCols{}
		cols.Col = append(cols.Col, col)
		ws.Cols = &cols
		return err
	}
	ws.Cols.Col = flatCols(col, ws.Cols.Col, func(fc, c xlsxCol) xlsxCol {
		fc.BestFit = c.BestFit
		fc.Collapsed = c.Collapsed
		fc.Hidden = c.Hidden
		fc.OutlineLevel = c.OutlineLevel
		fc.Phonetic = c.Phonetic
		fc.Style = c.Style
		return fc
	})
	return err
}

// flatCols provides a method for the column's operation functions to flatten
// and check the worksheet columns.
func flatCols(col xlsxCol, cols []xlsxCol, replacer func(fc, c xlsxCol) xlsxCol) []xlsxCol {
	fc := []xlsxCol{}
	for i := col.Min; i <= col.Max; i++ {
		c := deepcopy.Copy(col).(xlsxCol)
		c.Min, c.Max = i, i
		fc = append(fc, c)
	}
	inFlat := func(colID int, cols []xlsxCol) (int, bool) {
		for idx, c := range cols {
			if c.Max == colID && c.Min == colID {
				return idx, true
			}
		}
		return -1, false
	}
	for _, column := range cols {
		for i := column.Min; i <= column.Max; i++ {
			if idx, ok := inFlat(i, fc); ok {
				fc[idx] = replacer(fc[idx], column)
				continue
			}
			c := deepcopy.Copy(column).(xlsxCol)
			c.Min, c.Max = i, i
			fc = append(fc, c)
		}
	}
	return fc
}

// positionObjectPixels calculate the vertices that define the position of a
// graphical object within the worksheet in pixels.
//
//          +------------+------------+
//          |     A      |      B     |
//    +-----+------------+------------+
//    |     |(x1,y1)     |            |
//    |  1  |(A1)._______|______      |
//    |     |    |              |     |
//    |     |    |              |     |
//    +-----+----|    OBJECT    |-----+
//    |     |    |              |     |
//    |  2  |    |______________.     |
//    |     |            |        (B2)|
//    |     |            |     (x2,y2)|
//    +-----+------------+------------+
//
// Example of an object that covers some of the area from cell A1 to B2.
//
// Based on the width and height of the object we need to calculate 8 vars:
//
//    colStart, rowStart, colEnd, rowEnd, x1, y1, x2, y2.
//
// We also calculate the absolute x and y position of the top left vertex of
// the object. This is required for images.
//
// The width and height of the cells that the object occupies can be
// variable and have to be taken into account.
//
// The values of col_start and row_start are passed in from the calling
// function. The values of col_end and row_end are calculated by
// subtracting the width and height of the object from the width and
// height of the underlying cells.
//
//    colStart        # Col containing upper left corner of object.
//    x1              # Distance to left side of object.
//
//    rowStart        # Row containing top left corner of object.
//    y1              # Distance to top of object.
//
//    colEnd          # Col containing lower right corner of object.
//    x2              # Distance to right side of object.
//
//    rowEnd          # Row containing bottom right corner of object.
//    y2              # Distance to bottom of object.
//
//    width           # Width of object frame.
//    height          # Height of object frame.
//
func (f *File) positionObjectPixels(sheet string, col, row, x1, y1, width, height int) (int, int, int, int, int, int) {
	// Adjust start column for offsets that are greater than the col width.
	for x1 >= f.getColWidth(sheet, col) {
		x1 -= f.getColWidth(sheet, col)
		col++
	}

	// Adjust start row for offsets that are greater than the row height.
	for y1 >= f.getRowHeight(sheet, row) {
		y1 -= f.getRowHeight(sheet, row)
		row++
	}

	// Initialise end cell to the same as the start cell.
	colEnd := col
	rowEnd := row

	width += x1
	height += y1

	// Subtract the underlying cell widths to find end cell of the object.
	for width >= f.getColWidth(sheet, colEnd+1) {
		colEnd++
		width -= f.getColWidth(sheet, colEnd)
	}

	// Subtract the underlying cell heights to find end cell of the object.
	for height >= f.getRowHeight(sheet, rowEnd) {
		height -= f.getRowHeight(sheet, rowEnd)
		rowEnd++
	}

	// The end vertices are whatever is left from the width and height.
	x2 := width
	y2 := height
	return col, row, colEnd, rowEnd, x2, y2
}

// getColWidth provides a function to get column width in pixels by given
// sheet name and column index.
func (f *File) getColWidth(sheet string, col int) int {
	xlsx, _ := f.workSheetReader(sheet)
	if xlsx.Cols != nil {
		var width float64
		for _, v := range xlsx.Cols.Col {
			if v.Min <= col && col <= v.Max {
				width = v.Width
			}
		}
		if width != 0 {
			return int(convertColWidthToPixels(width))
		}
	}
	// Optimisation for when the column widths haven't changed.
	return int(defaultColWidthPixels)
}

// GetColWidth provides a function to get column width by given worksheet name
// and column index.
func (f *File) GetColWidth(sheet, col string) (float64, error) {
	colNum, err := ColumnNameToNumber(col)
	if err != nil {
		return defaultColWidthPixels, err
	}
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return defaultColWidthPixels, err
	}
	if ws.Cols != nil {
		var width float64
		for _, v := range ws.Cols.Col {
			if v.Min <= colNum && colNum <= v.Max {
				width = v.Width
			}
		}
		if width != 0 {
			return width, err
		}
	}
	// Optimisation for when the column widths haven't changed.
	return defaultColWidthPixels, err
}

// InsertCol provides a function to insert a new column before given column
// index. For example, create a new column before column C in Sheet1:
//
//    err := f.InsertCol("Sheet1", "C")
//
func (f *File) InsertCol(sheet, col string) error {
	num, err := ColumnNameToNumber(col)
	if err != nil {
		return err
	}
	return f.adjustHelper(sheet, columns, num, 1)
}

// RemoveCol provides a function to remove single column by given worksheet
// name and column index. For example, remove column C in Sheet1:
//
//    err := f.RemoveCol("Sheet1", "C")
//
// Use this method with caution, which will affect changes in references such
// as formulas, charts, and so on. If there is any referenced value of the
// worksheet, it will cause a file error when you open it. The excelize only
// partially updates these references currently.
func (f *File) RemoveCol(sheet, col string) error {
	num, err := ColumnNameToNumber(col)
	if err != nil {
		return err
	}

	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	for rowIdx := range ws.SheetData.Row {
		rowData := &ws.SheetData.Row[rowIdx]
		for colIdx := range rowData.C {
			colName, _, _ := SplitCellName(rowData.C[colIdx].R)
			if colName == col {
				rowData.C = append(rowData.C[:colIdx], rowData.C[colIdx+1:]...)[:len(rowData.C)-1]
				break
			}
		}
	}
	return f.adjustHelper(sheet, columns, num, -1)
}

// convertColWidthToPixels provieds function to convert the width of a cell
// from user's units to pixels. Excel rounds the column width to the nearest
// pixel. If the width hasn't been set by the user we use the default value.
// If the column is hidden it has a value of zero.
func convertColWidthToPixels(width float64) float64 {
	var padding float64 = 5
	var pixels float64
	var maxDigitWidth float64 = 7
	if width == 0 {
		return pixels
	}
	if width < 1 {
		pixels = (width * 12) + 0.5
		return math.Ceil(pixels)
	}
	pixels = (width*maxDigitWidth + 0.5) + padding
	return math.Ceil(pixels)
}
