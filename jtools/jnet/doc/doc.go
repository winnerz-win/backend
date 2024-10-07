package doc

//DocumentData :
type DocumentData struct {
	Href     string
	Size     string
	Weight   string
	Color    string
	Text     string
	onlyText bool
}

//DocStringList :
type DocStringList []DocumentData

//Count :
func (my DocStringList) Count() int {
	cnt := 0
	for _, v := range my {
		if v.onlyText == false {
			cnt++
		}
	} //for
	return cnt
}
