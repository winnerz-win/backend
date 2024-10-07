package mongo

type fStack []interface{}

func (my *fStack) Push(v interface{}) {
	*my = append(*my, v)
}
func (my *fStack) Pop() interface{} {
	idx := len(*my) - 1
	if idx < 0 {
		return nil
	}
	item := (*my)[idx]
	*my = (*my)[:idx]
	return item
}
func (my fStack) Count() int { return len(my) }

func (my fStack) Cmp(v interface{}) bool {
	if len(my) == 0 {
		return false
	}
	return my[len(my)-1] == v
}
