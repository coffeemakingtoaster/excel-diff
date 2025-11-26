package diff

import (
	"fmt"
	"slices"

	"github.com/xuri/excelize/v2"
)

type ExcelDiffFile struct {
	f         *excelize.File
	keys      map[string]*ExcelLine
	columnIds []int
}

func (edf ExcelDiffFile) getRowData(row []string, relevantColumnIndizes []int) (string, []string, error) {
	key := ""
	content := []string{}
	for i, val := range row {
		if slices.Contains(relevantColumnIndizes, i) {
			key += val
			content = append(content, val)
		}
	}
	return key, content, nil
}

func (edf *ExcelDiffFile) buildMap(relevantColumns []string) error {
	ids := edf.columnIds
	rows, err := edf.f.Rows(edf.f.GetSheetName(edf.f.GetActiveSheetIndex()))
	if err != nil {
		return err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}
		if len(ids) == 0 && len(ids) == 0 {
			for i, col := range row {
				if slices.Contains(relevantColumns, col) {
					ids = append(ids, i)
				}
			}
			if len(ids) != len(relevantColumns) {
				return fmt.Errorf("Not all columns are present")
			}
		}
		key, value, err := edf.getRowData(row, ids)
		if err != nil {
			return err
		}
		if _, ok := edf.keys[key]; !ok {
			edf.keys[key] = &ExcelLine{
				content: value,
				count:   0,
			}
		}
		edf.keys[key].count++
	}
	edf.columnIds = ids
	return nil
}

func (edf *ExcelDiffFile) tidyMap() {
	obsoleteKeys := []string{}
	for k, v := range edf.keys {
		if v.count == 0 {
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
			n := v.count
			that.keys[k].count -= this.keys[k].count
			this.keys[k].count -= n
		}
	}
	this.tidyMap()
	that.tidyMap()
}
