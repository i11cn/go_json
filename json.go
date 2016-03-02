package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

/*
用来测试的数据:

{
  "config": {"path": "./logs/", "format":"%T %L %N : %M"},
  "server": [2,
    "/var/www/html",
    {"host":"192.168.1.10", "port":10000},
    {"host":"192.168.1.11", "port":10000, "relay":{"host":"192.168.10.10", "port":20000}}
  ]
}
*/

type (
	Json struct {
		data json_value
	}
)

func NewJson() *Json {
	return &Json{&json_object{}}
}

func FromString(str string) (ret *Json, err error) {
	use := make(map[string]interface{})
	if err = json.Unmarshal([]byte(str), &use); err == nil {
		ret = FromMap(use)
	}
	return
}

func FromFile(str string) (*Json, error) {
	return NewJson(), nil
}

func FromMap(src map[string]interface{}) *Json {
	data := transform_from_map(src)
	return &Json{data}
}

func (j *Json) ToJson() (ret []byte, err error) {
	return json.Marshal(j.data)
}

func (j *Json) ToJson2() string {
	if d, e := j.ToJson(); e != nil {
		return ""
	} else {
		return string(d)
	}
}

func (j *Json) Get(path ...string) *Json {
	ret := j
	for _, k := range path {
		if ret = ret.get_child_by_key(k); ret == nil {
			return nil
		}
	}
	return ret
}

func (j *Json) Exist(path ...string) bool {
	return j.Get(path...) != nil
}

func (j *Json) Int() (ret int, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = int(reflect.ValueOf(v).Float())
		case uint, uint32, uint64:
			ret = int(reflect.ValueOf(v).Uint())
		case int, int32, int64:
			ret = int(reflect.ValueOf(v).Int())
		default:
			ret = 0
			ok = false
		}
	}
	return
}

func (j *Json) UInt() (ret uint, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = uint(reflect.ValueOf(v).Float())
		case uint, uint32, uint64:
			ret = uint(reflect.ValueOf(v).Uint())
		case int, int32, int64:
			ret = uint(reflect.ValueOf(v).Int())
		default:
			ret = 0
			ok = false
		}
	}
	return
}

func (j *Json) Int64() (ret int64, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = int64(reflect.ValueOf(v).Float())
		case uint, uint32, uint64:
			ret = int64(reflect.ValueOf(v).Uint())
		case int, int32, int64:
			ret = int64(reflect.ValueOf(v).Int())
		default:
			ret = 0
			ok = false
		}
	}
	return
}

func (j *Json) UInt64() (ret uint64, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = uint64(reflect.ValueOf(v).Float())
		case uint, uint32, uint64:
			ret = uint64(reflect.ValueOf(v).Uint())
		case int, int32, int64:
			ret = uint64(reflect.ValueOf(v).Int())
		default:
			ret = 0
			ok = false
		}
	}
	return
}

func (j *Json) Float() (ret float32, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = float32(reflect.ValueOf(v).Float())
		case uint, uint32, uint64:
			ret = float32(reflect.ValueOf(v).Uint())
		case int, int32, int64:
			ret = float32(reflect.ValueOf(v).Int())
		default:
			ret = 0
			ok = false
		}
	}
	return
}

func (j *Json) Float64() (ret float64, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = float64(reflect.ValueOf(v).Float())
		case uint, uint32, uint64:
			ret = float64(reflect.ValueOf(v).Uint())
		case int, int32, int64:
			ret = float64(reflect.ValueOf(v).Int())
		default:
			ret = 0
			ok = false
		}
	}
	return
}

func (j *Json) Bool() (ret bool, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case float32, float64:
			ret = reflect.ValueOf(v).Float() != 0
		case uint, uint32, uint64:
			ret = reflect.ValueOf(v).Uint() != 0
		case int, int32, int64:
			ret = reflect.ValueOf(v).Int() != 0
		case string:
			switch strings.ToUpper(v) {
			case "T", "Y", "TRUE", "YES", "1":
				ret = true
			case "F", "N", "FALSE", "NO", "0":
				ret = false
			default:
				ret = false
				ok = false
			}
		default:
			ret = false
			ok = false
		}
	}
	return
}

func (j *Json) Set(key string, value interface{}) *Json {
	switch v := j.data.(type) {
	case *json_object:
		v.set(key, value)
	case *json_array:
		v.set(key, value)
	default:
		j.data = create_json_array(v, create_json_object(key, value))
	}
	return j
}

func (j *Json) Append(key string, value interface{}) *Json {
	switch v := j.data.(type) {
	case *json_object:
		v.append(key, value)
	case *json_array:
		v.set_or_append(key, value)
	default:
		j.data = create_json_array(v, create_json_object(key, value))
	}
	return j
}

func (j *Json) SetValue(value interface{}) *Json {
	return j
}

func (j *Json) AppendValue(value interface{}) *Json {
	return j
}

func (j *Json) SetByPath(value interface{}, path ...string) *Json {
	l := len(path)
	switch l {
	case 0:
	case 1:
		j.Set(path[0], value)
	default:
		use := j
		fmt.Println(&use)
		for i := 0; use != nil && i < l-1; i++ {
			use = use.get_or_create_child(path[i])
			fmt.Println(&use)
		}
		if use != nil {
			use.Set(path[l-1], value)
		}
	}
	return j
}

func (j *Json) AppendByPath(value interface{}, path ...string) *Json {
	l := len(path)
	switch l {
	case 0:
	case 1:
		j.Append(path[0], value)
	default:
		use := j
		fmt.Println(&use)
		for i := 0; use != nil && i < l-1; i++ {
			use = use.get_or_create_child(path[i])
			fmt.Println(&use)
		}
		if use != nil {
			use.Append(path[l-1], value)
		}
	}
	return j
}

func (j *Json) Merge(j2 *Json) *Json {
	return j
}
