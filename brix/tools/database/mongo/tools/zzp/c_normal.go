package zzp

type normal struct {
	cmd    string
	params []string
	brunch brunch_data // [cmd]Result
	next   *normal
}

type brunch_data map[string][]string

func (my brunch_data) Callback(key string, f func([]string)) {
	if v, do := my[key]; do {
		if len(v) > 0 {
			f(v)
		}
	}
}

type Normaler interface {
	Cmd() string
	Params() []string
	Brunch() brunch_data
	Next() Normaler
	String() string
}

type normalize struct {
	*normal
	err error
}

func (my normalize) Error() error { return my.err }

func (my normalize) String() string {
	if my.err != nil {
		return my.err.Error()
	}
	return my.normal.String()
}

type Normalizer interface {
	Normaler
	Error() error
}

func (my normal) Cmd() string         { return my.cmd }
func (my normal) Params() []string    { return my.params }
func (my normal) Brunch() brunch_data { return my.brunch }
func (my normal) Next() Normaler      { return my.next }

func (my normal) void() interface{} {
	void := map[string]interface{}{}
	void["cmd"] = my.cmd
	if len(my.params) > 0 {
		void["iargs"] = my.params
	}
	if len(my.brunch) > 0 {
		void["brunch"] = my.brunch
	}
	if my.next != nil {
		void["next"] = my.next.void()
	}
	return void
}

func (my normal) String() string {
	return toJsonString(my.void())
}
func newNormal() *normal {
	return &normal{
		brunch: brunch_data{},
	}
}

func newNormalizer() *normalize {
	return &normalize{
		normal: newNormal(),
		err:    nil,
	}
}
