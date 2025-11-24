package diff

import (
	"fmt"
	"slices"

	"github.com/xuri/excelize/v2"
)

type ExcelLine struct {
	content  string
	position int
}

type ExcelDiffFile struct {
	f    *excelize.File
	keys map[string]int
}

func (edf ExcelDiffFile) getRowHash(row []string, relevantColumnIndizes []int) (string, error) {
	res := ""
	for i, val := range row {
		if slices.Contains(relevantColumnIndizes, i) {
			res += val
		}
	}
	return res, nil
}

func (edf *ExcelDiffFile) buildMap(relevantColumns []string) error {
	ids := []int{}
	rows, err := edf.f.Rows(edf.f.GetSheetName(edf.f.GetActiveSheetIndex()))
	if err != nil {
		return err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			for i, col := range row {
				if slices.Contains(relevantColumns, col) {
					ids = append(ids, i)
				}
			}
			if len(ids) != len(relevantColumns) {
				return fmt.Errorf("Not all columns are present")
			}
		}
		key, err := edf.getRowHash(row, ids)
		if err != nil {
			return err
		}
		if _, ok := edf.keys[key]; !ok {
			edf.keys[key] = 0
		}
		edf.keys[key]++

	}
	return nil
}

func (edf *ExcelDiffFile) tidyMap() {
	obsoleteKeys := []string{}
	for k, v := range edf.keys {
		if v == 0 {
			obsoleteKeys = append(obsoleteKeys, k)
		}
	}
	for _, k := range obsoleteKeys {
		delete(edf.keys, k)
	}
}

func (this *ExcelDiffFile) diff(that *ExcelDiffFile) {
	for k, v := range that.keys {
		if _, ok := this.keys[k]; ok {
			that.keys[k] -= this.keys[k]
			this.keys[k] -= v
		}
	}
	this.tidyMap()
	that.tidyMap()
}

type ExcelDiff struct {
	this            *ExcelDiffFile
	that            *ExcelDiffFile
	relevantColumns []string
	computed        bool
}

func loadExcelFile(path string) (*ExcelDiffFile, error) {
	fmt.Printf("Loading %s\n", path)
	edf := ExcelDiffFile{}
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	edf.f = f
	edf.keys = make(map[string]int)
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
	for k, v := range ed.that.keys {
		for range v {
			res = append(res, k)
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
	for k, v := range ed.this.keys {
		for range v {
			res = append(res, k)
		}
	}
	return res
}
