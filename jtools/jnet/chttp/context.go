package chttp

import (
	"context"
	"strings"
)

type _Keyword string

const (
	_CLASSIC_NAME    = _Keyword("classic_name")
	_PORT            = _Keyword("port")
	_HEADER          = _Keyword("header")
	_KEEP_ALIVE      = _Keyword("keepAlive")
	_SERVE_FILE_PATH = _Keyword("serveFilePath")
	_SSL             = _Keyword("ssl")
	_STARTER         = _Keyword("starter")
	_CONSOLE_SKIP    = _Keyword("console_skip")
)

func isAllowKeyword(n string) bool {
	switch _Keyword(n) {
	case _CLASSIC_NAME, _PORT,
		_HEADER, _KEEP_ALIVE,
		_SERVE_FILE_PATH, _SSL,
		_STARTER:
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////

type iContext struct {
	context.Context
}

func Context() *iContext {
	my := &iContext{context.Background()}
	my.SetName("classic").
		SetPort(8080).
		SetHeaders("Access-Control-Allow-Origin").
		SetServeFilePath(DefaultAssetsName).
		SetConsoleSkip(false).
		SetKeepAlive(false).
		SetSSL("", "").
		SetStarter()
	return my
}

func (my *iContext) set(key string, val interface{}) {
	my.Context = context.WithValue(my.Context, key, val)
}
func (my *iContext) _set(key _Keyword, val interface{}) {
	my.set(string(key), val)
}
func (my *iContext) get(key string) interface{} {
	return my.Context.Value(key)
}
func (my *iContext) _get(key _Keyword) interface{} {
	return my.get(string(key))
}

func (my *iContext) Set(key string, val interface{}) {
	if !isAllowKeyword(key) {
		return
	}
	my.set(key, val)
}

func (my *iContext) Get(key string) interface{} {
	return my.get(key)
}

func (my *iContext) SetName(name string) *iContext {
	my._set(_CLASSIC_NAME, name)
	return my
}
func (my iContext) Name() string { return my._get(_CLASSIC_NAME).(string) }

func (my *iContext) SetHeaders(headers ...string) *iContext {
	list := []string{}
	val := my._get(_HEADER)
	if val != nil {
		list = val.([]string)
	}
	list = append(list, headers...)
	my._set(_HEADER, list)
	return my
}
func (my iContext) Headers() []string {
	return my._get(_HEADER).([]string)
}

func (my *iContext) SetConsoleSkip(b bool) *iContext {
	my._set(_CONSOLE_SKIP, b)
	return my
}
func (my iContext) ConsoleSkip() bool {
	return my._get(_CONSOLE_SKIP).(bool)
}

func (my *iContext) SetPort(port int) *iContext {
	my._set(_PORT, port)
	return my
}
func (my iContext) Port() int {
	return my._get(_PORT).(int)
}

func (my *iContext) SetKeepAlive(b bool) *iContext {
	my._set(_KEEP_ALIVE, b)
	return my
}
func (my iContext) KeepAlive() bool {
	return my._get(_KEEP_ALIVE).(bool)
}

func (my *iContext) SetServeFilePath(serveFilePath string) *iContext {
	my._set(_SERVE_FILE_PATH, serveFilePath)
	return my
}
func (my iContext) ServeFilePath() string {
	return my._get(_SERVE_FILE_PATH).(string)
}

func (my *iContext) SetSSL(certFile, keyFile string) *iContext {
	certFile = strings.TrimSpace(certFile)
	keyFile = strings.TrimSpace(keyFile)
	if certFile == "" || keyFile == "" {
		if my._get(_SSL) == nil {
			my._set(_SSL, []string{})
		}
		return my
	}
	my._set(_SSL, []string{certFile, keyFile})
	return my
}

func (my iContext) SSL() []string {
	return my._get(_SSL).([]string)
}

func (my *iContext) SetStarter(starters ...Starter) *iContext {
	if len(starters) == 0 {
		if my._get(_STARTER) == nil {
			my._set(_STARTER, StarterList{})
		}
		return my
	}
	list := StarterList{}
	v := my._get(_STARTER)
	if v != nil {
		list = v.(StarterList)
	}
	list = append(list, starters...)
	my._set(_STARTER, list)
	return my
}
func (my *iContext) StarterList() StarterList {
	return my._get(_STARTER).(StarterList)
}

////////////////////////////////////////////////////////////////////////

func (my iContext) getMap() map[_Keyword]interface{} {
	v := map[_Keyword]interface{}{}

	return v
}

func (my iContext) String() string {
	return ToJsonString(my.getMap())
}

////////////////////////////////////////////////////////////////////////
