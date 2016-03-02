package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
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
	Options struct {
		name  string
		child interface{}
	}
)

func NewOptions() *Options {
	return &Options{child: map[string]*Options{}}
}

func (o *Options) Json() (ret string, err error) {
	var d []byte
	if d, err = json.Marshal(o); err == nil {
		ret = string(d)
	}
	return
}

func (o *Options) MarshalJSON() (ret []byte, err error) {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("{")
	if err = o.marshal_json(buf); err == nil {
		buf.WriteString("}")
		ret = buf.Bytes()
	}
	return
}

func (o *Options) Get(path ...string) interface{} {
	if o == nil || len(path) == 0 {
		return nil
	}
	current := o
	for _, k := range path {
		if current = current.get_child_by_key(k, false); current == nil {
			return nil
		}
	}
	return current
}

func (o *Options) Exist(path ...string) bool {
	return o.Get(path...) != nil
}

func (o *Options) Key() string {
	return o.name
}

func (o *Options) String() string {
	if use, exist := o.get_value(); exist {
		switch v := use.(type) {
		case string:
			return v
		default:
			return fmt.Sprint(v)
		}
	}
	return ""
}

func (o *Options) Int() (ret int, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) UInt() (ret uint, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) Int64() (ret int64, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) UInt64() (ret uint64, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) Float() (ret float32, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) Float64() (ret float64, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) Bool() (ret bool, ok bool) {
	ok = true
	if use, exist := o.get_value(); exist {
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

func (o *Options) ParseJsonFile(path string) error {
	return nil
}

func (o *Options) Set(key string, value interface{}) *Options {
	o.name = key
	return o
}

func (o *Options) Replace(key string, value interface{}) *Options {
	return o
}

func (o *Options) Append(key string, value interface{}) *Options {
	return o
}

func (o *Options) SetPath(value interface{}, path ...string) *Options {
	return o
}

func (o *Options) ReplacePath(value interface{}, path ...string) *Options {
	return o
}

func (o *Options) AppendPath(value interface{}, path ...string) *Options {
	return o
}

func (o *Options) Merge(o2 *Options) *Options {
	return o
}

func (o *Options) SetJson(j *simplejson.Json) *Options {
	return o
}

func (o *Options) ParseJson(js string) error {
	return o.parse_json_string(js)
}

func (o *Options) ParseIniFile(path string) error {
	return nil
}

func (o *Options) ParseCommand() *Options {
	return o
}
