// Copyright 2016 - 2020 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to
// and read from XLSX files. Support reads and writes XLSX file generated by
// Microsoft Excel™ 2007 and later. Support save file without losing original
// charts of XLSX. This library needs Go version 1.10 or later.

package excelize

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddComments(t *testing.T) {
	f, err := prepareTestBook1()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	s := strings.Repeat("c", 32768)
	assert.NoError(t, f.AddComment("Sheet1", "A30", `{"author":"`+s+`","text":"`+s+`"}`))
	assert.NoError(t, f.AddComment("Sheet2", "B7", `{"author":"Excelize: ","text":"This is a comment."}`))

	// Test add comment on not exists worksheet.
	assert.EqualError(t, f.AddComment("SheetN", "B7", `{"author":"Excelize: ","text":"This is a comment."}`), "sheet SheetN is not exist")
	// Test add comment on with illegal cell coordinates
	assert.EqualError(t, f.AddComment("Sheet1", "A", `{"author":"Excelize: ","text":"This is a comment."}`), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	if assert.NoError(t, f.SaveAs(filepath.Join("test", "TestAddComments.xlsx"))) {
		assert.Len(t, f.GetComments(), 2)
	}

	f.Comments["xl/comments2.xml"] = nil
	f.XLSX["xl/comments2.xml"] = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><comments xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><authors><author>Excelize: </author></authors><commentList><comment ref="B7" authorId="0"><text><t>Excelize: </t></text></comment></commentList></comments>`)
	comments := f.GetComments()
	assert.EqualValues(t, 2, len(comments["Sheet1"]))
	assert.EqualValues(t, 1, len(comments["Sheet2"]))
}

func TestDecodeVMLDrawingReader(t *testing.T) {
	f := NewFile()
	path := "xl/drawings/vmlDrawing1.xml"
	f.XLSX[path] = MacintoshCyrillicCharset
	f.decodeVMLDrawingReader(path)
}

func TestCommentsReader(t *testing.T) {
	f := NewFile()
	path := "xl/comments1.xml"
	f.XLSX[path] = MacintoshCyrillicCharset
	f.commentsReader(path)
}

func TestCountComments(t *testing.T) {
	f := NewFile()
	f.Comments["xl/comments1.xml"] = nil
	assert.Equal(t, f.countComments(), 1)
}
