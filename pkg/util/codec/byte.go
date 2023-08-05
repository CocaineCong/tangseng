package codec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/spf13/cast"
)

// BinaryWrite 将所有的类型 转成byte buffer类型，易于存储// TODO change
func BinaryWrite(v any) (buf *bytes.Buffer, err error) {
	if v == nil {
		err = errors.New("BinaryWrite the value is nil")
		return
	}
	buf = new(bytes.Buffer)

	switch reflect.Indirect(reflect.ValueOf(v)).Kind() { // TODO:反射很影响性能，后续看看怎么优化
	case reflect.Int64, reflect.Int32, reflect.Int:
		buf.Write([]byte(cast.ToString(v)))
	case reflect.String:
		buf.Write([]byte(v.(string)))
	case reflect.Slice, reflect.Array, reflect.Struct:
		value, errx := json.Marshal(v)
		if errx != nil {
			err = errx
			return
		}
		buf.Write(value)
	}

	return
}

// GobWrite 将所有的类型 转成 bytes.Buffer 类型，易于存储// TODO change
func GobWrite(v any) (buf *bytes.Buffer, err error) {
	if v == nil {
		err = errors.New("BinaryWrite the value is nil")
		return
	}
	buf = new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(v); err != nil {
		return
	}

	return
}
