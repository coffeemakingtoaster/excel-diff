package diff

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelDiff struct {
	this            *ExcelDiffFile
	that            *ExcelDiffFile
	relevantColumns []string
	computed        bool
}

func loadExcelFile(path string) (*ExcelDiffFile, error) {
	edf := ExcelDiffFile{}
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	edf.f = f
	edf.keys = make(map[string]*ExcelLine)
	return &edf, nil
}

func LoadExcelDiff(thisPath, thatPath string, columns []string) (*ExcelDiff, error) {
	this, err := loadExcelFile(thisPath)
	if err != nil {
		return nil, err
	}
	that, err := loadExcelFile(thatPath)
	if err != nil {
		return nil, err
	}

	return &ExcelDiff{this, that, columns, false}, nil
}

func (ed *ExcelDiff) compute() {
	ed.this.buildMap(ed.relevantColumns)
	ed.that.buildMap(ed.relevantColumns)
	ed.this.diff(ed.that)
	ed.computed = true
}

// Get lines that that contains and this doesn't
func (ed *ExcelDiff) GetAddedLines() []string {
	if !ed.computed {
		ed.compute()
	}
	res := []string{}
	for _, v := range ed.that.keys {
		for range v.count {
			res = append(res, strings.Join(v.content, SEP))
		}
	}
	return res
}

// Get lines that this contains and that doesn't
func (ed *ExcelDiff) GetRemovedLines() []string {
	if !ed.computed {
		ed.compute()
	}
	res := []string{}
	for _, v := range ed.this.keys {
		for range v.count {
			res = append(res, strings.Join(v.content, SEP))
		}
	}
	return res
}
