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

import "encoding/xml"

// xlsxProperties specifies to an OOXML document properties such as the
// template used, the number of pages and words, and the application name and
// version.
type xlsxProperties struct {
	XMLName              xml.Name `xml:"http://schemas.openxmlformats.org/officeDocument/2006/extended-properties Properties"`
	Template             string
	Manager              string
	Company              string
	Pages                int
	Words                int
	Characters           int
	PresentationFormat   string
	Lines                int
	Paragraphs           int
	Slides               int
	Notes                int
	TotalTime            int
	HiddenSlides         int
	MMClips              int
	ScaleCrop            bool
	HeadingPairs         *xlsxVectorVariant
	TitlesOfParts        *xlsxVectorLpstr
	LinksUpToDate        bool
	CharactersWithSpaces int
	SharedDoc            bool
	HyperlinkBase        string
	HLinks               *xlsxVectorVariant
	HyperlinksChanged    bool
	DigSig               *xlsxDigSig
	Application          string
	AppVersion           string
	DocSecurity          int
}

// xlsxVectorVariant specifies the set of hyperlinks that were in this
// document when last saved.
type xlsxVectorVariant struct {
	Content string `xml:",innerxml"`
}

type xlsxVectorLpstr struct {
	Content string `xml:",innerxml"`
}

// xlsxDigSig contains the signature of a digitally signed document.
type xlsxDigSig struct {
	Content string `xml:",innerxml"`
}
