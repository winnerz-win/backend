package lvdb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LEVELDB struct {
	path string //db_path
	*leveldb.DB
}

func New(filepath string, isReadOnly ...bool) *LEVELDB {

	var option *opt.Options = nil
	if len(isReadOnly) > 0 && isReadOnly[0] {
		option = &opt.Options{
			ReadOnly: true,
		}
	}

	db, err := leveldb.OpenFile(filepath, option)
	if err != nil {
		fmt.Println("lvdb.New :", err)
		return nil
	}

	return &LEVELDB{
		path: filepath,
		DB:   db,
	}
}

func (my *LEVELDB) PATH() string { return my.path }

func (my *LEVELDB) IsEmpty() bool {
	iter := my.NewIterator(nil, nil)
	defer iter.Release()
	do := iter.First() || iter.Last()
	return do
}

func (my *LEVELDB) ACTION(f func(lv *LEVELDB)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("lvdb.Action : %v", e)
		}
		//my.Close()
	}()

	f(my)

	return nil
}

type _Integer interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
	///~int | ~uint |
}
type _Double interface {
	~float32 | ~float64
}

type _NumberType interface {
	_Integer | _Double | ~bool
}

func BytesToNumber[T _NumberType](v []byte) (T, error) {
	buf := new(bytes.Buffer)
	buf.Write(v)

	var data T
	if err := binary.Read(buf, binary.BigEndian, &data); err != nil {
		return data, fmt.Errorf("BytesToNumber : %v", err)
	}
	return data, nil
}
func BytesToNumber1[T _NumberType](v []byte) T {
	re, err := BytesToNumber[T](v)
	if err != nil {
		fmt.Println("lvdb.BytesToNumber1 : ", err)
	}
	return re
}

func NumberToBytes[T _NumberType](v T) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, v); err != nil {
		panic(string(debug.Stack()))
	}
	return buf.Bytes()
}

type _JsonBytes interface {
	Bytes() []byte
}

func ToBytes(p interface{}) []byte {

	switch v := p.(type) {
	case []byte:
		return v

	case string:
		return []byte(v)

	case int8:
		return NumberToBytes(v)
	case uint8:
		return NumberToBytes(v)
	case int16:
		return NumberToBytes(v)
	case uint16:
		return NumberToBytes(v)
	case int32:
		return NumberToBytes(v)
	case uint32:
		return NumberToBytes(v)

	case int64:
		return NumberToBytes(v)
	case uint64:
		return NumberToBytes(v)
	case int:
		return NumberToBytes(int64(v))
	case uint:
		return NumberToBytes(uint64(v))

	case float32:
		return NumberToBytes(v)
	case float64:
		return NumberToBytes(v)

	case bool:
		return NumberToBytes(v)

	case _JsonBytes:
		return v.Bytes()

	default:
		buf, err := json.Marshal(p)
		if err != nil {
			panic(string(debug.Stack()))
		}
		return buf
	} //switch
}

// Put(key []byte, value []byte, wo *opt.WriteOptions) error
func (my *LEVELDB) PUT(key, value interface{}, wo *opt.WriteOptions) error {
	return my.Put(
		ToBytes(key),
		ToBytes(value),
		wo,
	)
}

// Delete(key []byte, wo *opt.WriteOptions) error
func (my *LEVELDB) DELETE(key interface{}, wo *opt.WriteOptions) error {
	return my.Delete(
		ToBytes(key),
		wo,
	)
}

// Has(key []byte, ro *opt.ReadOptions) (ret bool, err error)
func (my *LEVELDB) HAS(key interface{}, ro *opt.ReadOptions) bool {
	do, err := my.Has(
		ToBytes(key),
		ro,
	)
	if err != nil {
		return false
	}
	return do
}

// Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
func (my *LEVELDB) GET(key interface{}, ro *opt.ReadOptions) (value []byte, err error) {
	return my.Get(
		ToBytes(key),
		ro,
	)
}

func GetJson[T any](db *LEVELDB, key interface{}, ro *opt.ReadOptions) (T, error) {
	var result T
	v, err := db.GET(key, ro)
	if err != nil {
		return result, err
	}
	json.Unmarshal(v, &result)
	return result, err
}
