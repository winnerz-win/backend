package excel

import "github.com/tealeg/xlsx"

type Row struct {
	row *xlsx.Row
}
type WCell struct {
	cell *xlsx.Cell
}

func (my *Row) R() *xlsx.Row { return my.row }

func (my *Row) Cell(index ...int) *WCell {
	var c *xlsx.Cell

	if len(index) > 0 {
		if index[0] <= 0 {
			c = my.row.AddCell()
		} else {
			count := index[0] + 1
			for i := count - 1; i >= 0; i-- {
				c = my.row.AddCell()
			}
		}
	} else {
		c = my.row.AddCell()
	}
	return &WCell{c}
}

func (my *Sheet) EmptyRow() {
	defer my.f.mu.Unlock()
	my.f.mu.Lock()
	my.AddRow()
}

func (my *Sheet) AppendRow(callback func(row *Row)) {
	defer my.f.mu.Unlock()
	my.f.mu.Lock()
	row := my.AddRow()
	callback(&Row{row})
}

func (my *WCell) C() *xlsx.Cell {
	return my.cell
}

func (my *WCell) SetValue(v interface{}) {
	my.cell.SetValue(v)
}

func (my *WCell) SetFontColor(c string) {
	my.cell.GetStyle().Font.Color = c
}
func (my *WCell) SetFontBold(b bool) {
	my.cell.GetStyle().Font.Bold = b
}
func (my *WCell) SetFontItalic(b bool) {
	my.cell.GetStyle().Font.Italic = b
}
func (my *WCell) SetFontUnderline(b bool) {
	my.cell.GetStyle().Font.Underline = b
}

func (my *WCell) SetFontSize(size int) {
	my.cell.GetStyle().Font.Size = size
}

func (my *WCell) SetFontName(name string) {
	my.cell.GetStyle().Font.Name = name
}

func (my *WCell) SetFontCallback(callback func(*xlsx.Font)) {
	callback(&my.cell.GetStyle().Font)
}
