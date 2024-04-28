package structutil

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"reflect"
)

func Struct2String(v any) string {
	bs, _ := json.Marshal(v)
	buf := new(bytes.Buffer)
	_ = json.Indent(buf, bs, "", "    ")
	return buf.String()
}

// Struct2Map convert struct to map v must be a pointer
func Struct2Map[T any](v T) map[string]any {
	res := make(map[string]any)

	elem := reflect.ValueOf(v).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		k := relType.Field(i).Name
		if k == "BaseModel" || k == "ID" {
			continue
		}
		res[k] = elem.Field(i).Interface()
	}

	return res
}

func Struct2Bytes[T any](v T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

func Bytes2Struct[T any](data []byte) (T, error) {
	var v T
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&v)
	return v, err
}

func BytesArray2Struct[T any](data [][]byte) ([]T, error) {
	res := make([]T, 0)
	for _, b := range data {
		v, err := Bytes2Struct[T](b)
		if err != nil {
			logrus.Errorf("BytesArray2Struct error: %+v", err)
			return res, err
		}
		res = append(res, v)
	}
	return res, nil
}

func AnyIsNil(v any) bool {
	switch v.(type) {
	case string:
		return len(v.(string)) == 0
	case int:
		return v.(int) == 0
	case int8:
		return v.(int8) == 0
	case int16:
		return v.(int16) == 0
	case int32:
		return v.(int32) == 0
	case int64:
		return v.(int64) == 0
	case uint:
		return v.(uint) == 0
	case uint8:
		return v.(uint8) == 0
	case uint16:
		return v.(uint16) == 0
	case uint32:
		return v.(uint32) == 0
	case uint64:
		return v.(uint64) == 0
	case float32:
		return v.(float32) == 0
	case float64:
		return v.(float64) == 0
	case bool:
		return !v.(bool)
	default:
		return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}

func ValueOfPtr[T string | int | int64 | int32](ptr *T, defaultVal T) T {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
