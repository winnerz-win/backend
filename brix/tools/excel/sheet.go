package excel

import "github.com/tealeg/xlsx"

type Sheet struct {
	f *Excel

	*xlsx.Sheet
	headerRow *xlsx.Row
}

func (my *Excel) GetSheet(name string) *Sheet {
	my.mu.RLock()
	if s, do := my.sh[name]; do {
		defer my.mu.RUnlock()
		return s
	}
	my.mu.RUnlock()

	/////////////////////////////////////////////

	my.mu.Lock()
	defer my.mu.Unlock()
	sheet, _ := my.File.AddSheet(name)
	my.sh[name] = &Sheet{
		f:     my,
		Sheet: sheet,
	}
	return my.sh[name]
}

func (my *Sheet) Headers(color string, names ...interface{}) {
	defer my.f.mu.Unlock()
	my.f.mu.Lock()

	if my.headerRow == nil {
		my.headerRow = my.AddRow()
	}
	for _, name := range names {
		cell := my.headerRow.AddCell()
		cell.SetValue(name)
		if color != "" {
			cell.GetStyle().Font.Color = color
		}
	} //for
}

func (my *Sheet) Cells(color string, names ...interface{}) *Row {
	defer my.f.mu.Unlock()
	my.f.mu.Lock()

	row := my.AddRow()
	for _, name := range names {
		cell := row.AddCell()
		cell.SetValue(name)
		if color != "" {
			cell.GetStyle().Font.Color = color
		}
		//cell.GetStyle().Font.Bold = true

	} //for
	return &Row{row}
}

type ISheet *xlsx.Sheet

func (my *Sheet) Action(f func(sheet ISheet)) {
	defer my.f.mu.Unlock()
	my.f.mu.Lock()

	f(my.Sheet)

}
