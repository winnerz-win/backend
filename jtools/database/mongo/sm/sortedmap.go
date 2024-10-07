package sm

import "sort"

//SortedML :
type SortedML struct {
	Keys   []string
	Values []interface{}
	KeyLen int
}

func (my SortedML) Len() int {
	return len(my.Keys)
}
func (my SortedML) Swap(i, j int) {
	my.Keys[i], my.Keys[j] = my.Keys[j], my.Keys[i]
	my.Values[i], my.Values[j] = my.Values[j], my.Values[i]
}
func (my SortedML) Less(i, j int) bool {
	if my.Keys[i] < my.Keys[j] {
		return true
	}
	return false
}

//MakeSortedML :
func MakeSortedML(src map[string]interface{}) SortedML {
	sm := SortedML{}
	for k, v := range src {
		if sm.KeyLen < len(k) {
			sm.KeyLen = len(k)
		}
		sm.Keys = append(sm.Keys, k)
		sm.Values = append(sm.Values, v)
	} //for
	sort.Sort(sm)
	return sm
}

//Foreach :
func (my SortedML) Foreach(fc func(key string, v interface{})) {
	size := len(my.Keys)
	for i := 0; i < size; i++ {
		fc(my.Keys[i], my.Values[i])
	}
}
