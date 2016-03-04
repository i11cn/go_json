/*
这个json包在处理json数据时，不需要和预先定义的struct对应，
而是将数据映射成了map和数组，通过map和数组直接操作数据

使用方法如下：
	import "github.com/i11cn/go_json"
	...

	j = json.NewJson() // 此处会创建一个没有数据的空Json结构
	j, err = json.FromString(`{"test":10, "other":"字符串"}`) // 此处使用给定的json字符串构造一个Json结构
	j, err = json.FromFile("config.json") // 此处从文件config.json中读取内容，构造一个Json结构
	...

	if j.Exist("test", "sub1") { // 检查test子节点下的sub1子节点是否存在
		sub := j.Get("test", "sub1") // 获取test节点下的sub1节点下的内容，并且构造成Json结构返回
		sub.Int() // 将该节点中的内容
	}

*/
package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
		data interface{}
	}
)

func NewJson() *Json {
	return &Json{make(map[string]interface{})}
}

func FromString(str string) (ret *Json, err error) {
	var use interface{}
	if err = json.Unmarshal([]byte(str), &use); err == nil {
		ret = FromObject(use)
	}
	return
}

func FromFile(str string) (ret *Json, err error) {
	var file *os.File
	if file, err = os.Open(str); err != nil {
		return
	}
	var data []byte
	if data, err = ioutil.ReadAll(file); err != nil {
		return
	}
	var use interface{}
	if err = json.Unmarshal(data, &use); err != nil {
		return
	}
	return FromObject(use), nil
}

func FromMap(src map[string]interface{}) *Json {
	data := transform_from_map(src)
	return &Json{data}
}

func FromObject(src interface{}) *Json {
	var data interface{}
	switch s := src.(type) {
	case []interface{}:
		data = transform_from_array(s)
	case map[string]interface{}:
		data = transform_from_map(s)
	default:
		data = create_json_array(s)
	}
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

func (j *Json) String() (ret string, ok bool) {
	ok = true
	if use, exist := j.get_value(); exist {
		switch v := use.(type) {
		case *[]interface{}:
			ok = false
		case map[string]interface{}:
			ok = false
		case string:
			ret = v
		default:
			ret = fmt.Sprint(v)
		}
	}
	return
}

func (j *Json) Array() []Json {
	ret := []Json{}
	if arr, ok := j.data.(*[]interface{}); ok {
		for _, d := range *arr {
			switch d.(type) {
			case *[]interface{}:
				ret = append(ret, Json{d})
			case map[string]interface{}:
				ret = append(ret, Json{d})
			default:
				ret = append(ret, Json{create_json_array(d)})
			}
		}
	}
	return ret
}

func (j *Json) Set(key string, value interface{}) *Json {
	switch v := j.data.(type) {
	case map[string]interface{}:
		v[key] = value
	case *[]interface{}:
		*v = set_array(*v, key, value)
	default:
		j.data = create_json_array(v, create_json_object(key, value))
	}
	return j
}

func (j *Json) Append(key string, value interface{}) *Json {
	switch v := j.data.(type) {
	case map[string]interface{}:
		v = append_object(v, key, value)
	case *[]interface{}:
		*v = set_or_append_array(*v, key, value)
	default:
		j.data = create_json_array(v, create_json_object(key, value))
	}
	return j
}

func (j *Json) AppendValue(value interface{}) *Json {
	switch d := j.data.(type) {
	case *[]interface{}:
		*d = append_array(*d, value)
	default:
	}
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
		for i := 0; use != nil && i < l-1; i++ {
			use = use.get_or_create_child(path[i])
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
		for i := 0; use != nil && i < l-1; i++ {
			use = use.get_or_create_child(path[i])
		}
		if use != nil {
			use.Append(path[l-1], value)
		}
	}
	return j
}

func (j *Json) IsObject() (ret bool) {
	_, ret = j.data.(map[string]interface{})
	return
}

func (j *Json) IsArray() (ret bool) {
	_, ret = j.data.([]interface{})
	return ret
}

func (j *Json) IsData() bool {
	return !j.IsObject() && !j.IsArray()
}

func (j *Json) Merge(j2 *Json) *Json {
	return j
}
