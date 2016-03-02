package option

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type (
	json_value  interface{}
	json_string string
	json_number struct {
		value interface{}
	}
	json_bool bool
	json_null interface{}

	json_array  []json_value
	json_object map[string]json_value

	Json struct {
		data json_value
	}
	Options2 Json
)

func NewOptions2() *Options2 {
	j := &Json{json_object{}}
	return (*Options2)(j)
}

func (o *Options2) Json() (ret string, err error) {
	if d, e := json.Marshal(o.data); e != nil {
		err = e
	} else {
		ret = string(d)
	}
	return
}

func (o *Options2) Json2() string {
	if d, e := json.Marshal(o.data); e != nil {
		return ""
	} else {
		return string(d)
	}
}

func (a *json_array) append(value ...interface{}) *json_array {
	for _, v := range value {
		*a = append(*a, v)
	}
	return a
}

func create_json_array(value ...interface{}) json_array {
	ret := json_array{}
	for _, v := range value {
		ret = append(ret, v)
	}
	return ret
}

func create_json_object(key string, value interface{}) json_object {
	ret := map[string]json_value{}
	ret[key] = value
	return json_object(ret)
}

func (src json_object) get_child_by_key(key string, create bool) json_value {
	obj := map[string]json_value(src)
	if data, exist := obj[key]; exist {
		return data
	} else if create {
		ret := create_json_object(key, json_object{})
		obj[key] = ret
		return ret
	}
	return nil
}

func (src json_array) get_child_by_key(key string, create bool) json_value {
	arr := ([]json_value)(src)
	var use json_object = nil
	for _, c := range arr {
		if v, ok := c.(json_object); ok {
			if use != nil {
				return nil
			}
			use = v
		}
	}
	return use.get_child_by_key(key, create)
}

func (o *Options2) Get(key string) *Options2 {
	j := (*Json)(o)
	switch data := j.data.(type) {
	case json_object:
		use := map[string]json_value(data)
		if v, exist := use[key]; exist {
			return &Options2{v}
		} else {
			return nil
		}
	case json_array:
		fmt.Println("是个数组，处理逻辑还没想清楚，暂时不处理...")
		return nil
	default:
		return nil
	}
}

func (o *Options2) GetOrCreate(key string) *Options2 {
	j := (*Json)(o)
	switch data := j.data.(type) {
	case json_object:
		fmt.Println("json_object -- ")
		use := map[string]json_value(data)
		if v, exist := use[key]; exist {
			return &Options2{v}
		} else {
			use[key] = json_object{}
			return &Options2{use[key]}
		}
	case json_array:
		fmt.Println("json_array -- ")
		j := data.get_child_by_key(key, true)
		return &Options2{j}
	default:
		fmt.Println("others -- ")
		ret := Options2{create_json_object(key, json_object{})}
		j.data = create_json_array(data, ret.data)
		return &ret
	}
}

func (o *Options2) set(value interface{}) *Options2 {
	j := (*Json)(o)
	if value == nil {
		value = json_null(value)
	}
	switch data := j.data.(type) {
	case json_array:
		data.append(value)
	case json_object:
		use := map[string]json_value(data)
		if len(use) == 0 {
			j.data = value
		} else {
			j.data = create_json_array(data, value)
		}
	case nil:
		j.data = value
	default:
		j.data = create_json_array(data, value)
	}
	return o
}

func (o *Options2) replace(value interface{}) *Options2 {
	j := (*Json)(o)
	if value == nil {
		value = json_null(value)
	}
	j.data = value
	return o
}

func (o *Options2) Set(key string, value interface{}) *Options2 {
	use := o.GetOrCreate(key)
	fmt.Println("use = ", use.Json2())
	fmt.Println("o = ", o.Json2())
	use.set(value)
	fmt.Println("use = ", use.Json2())
	fmt.Println("o = ", o.Json2())
	//j := (*Json)(o)
	//if value == nil {
	//	value = json_null(value)
	//}
	//switch data := j.data.(type) {
	//case json_object:
	//	use := map[string]json_value(data)
	//	if v, exist := use[key]; exist {
	//		switch c := v.(type) {
	//		case json_array:
	//			use[key] = c.append(value)
	//		default:
	//			use[key] = create_json_array(value, v)
	//		}
	//	} else {
	//		use[key] = value
	//	}
	//case json_array:
	//	fmt.Println("是个数组，处理逻辑还没想清楚，暂时不处理...")
	//case json_null:
	//	fmt.Println("是个空，保留，暂时不处理...")
	//case nil:
	//	v := json_object{}
	//	use := map[string]json_value(v)
	//	use[key] = value
	//	o.data = v
	//default:
	//	fmt.Println("是个其他数据，处理...")
	//	o.data = create_json_array(data, create_json_object(key, value))
	//}
	return o
}

func (o *Options2) Replace(key string, value interface{}) *Options2 {
	j := (*Json)(o)
	switch data := j.data.(type) {
	case json_object:
		fmt.Println("是个对象，处理...")
		use := map[string]json_value(data)
		use[key] = value
	case json_array:
		fmt.Println("是个数组，暂时不处理...")
	case nil:
		fmt.Println("是个空，处理...")
		v := json_object{}
		use := map[string]json_value(v)
		use[key] = value
		o.data = v
	default:
		fmt.Println("是个其他数据，处理...")
		v := json_object{}
		use := map[string]json_value(v)
		use[key] = value
		d := json_array{}
		use2 := []json_value(d)
		use2 = append(use2, v)
		use2 = append(use2, data)
		o.data = use2
	}
	return o
}

func dump_array(indent string, a []interface{}) {
	for _, v := range a {
		switch use := v.(type) {
		case map[string]interface{}:
			fmt.Printf("%s(map[string]interface{})  =\r\n", indent)
			dump_map(indent+"  ", &use)
		case []interface{}:
			fmt.Print("%s([]interface{}) =\r\n", indent)
			dump_array(indent+"  ", use)
		default:
			fmt.Printf("%s(%T) = %v\r\n", indent, use, use)
		}
	}
}

func dump_map(indent string, m *map[string]interface{}) {
	for k, v := range *m {
		switch use := v.(type) {
		case map[string]interface{}:
			fmt.Printf("%s(map[string]interface{}) - %s =\r\n", indent, k)
			dump_map(indent+"  ", &use)
		case []interface{}:
			fmt.Printf("%s([]interface{}) - %s =", indent, k)
			dump_array(indent+"  ", use)
		default:
			fmt.Printf("%s(%T) - %s = %v\r\n", indent, use, k, use)
		}
	}
}

func TestJson() {
	str := `{
  "config": {"path": "./logs/", "format":"%T %L %N : %M"},
  "server": [2,
    "/var/www/html",
	{"host":"192.168.1.10", "port":10000, "enable":false},
	{"host":"192.168.1.11", "port":10000, "ha":null, "relay":{"host":"192.168.10.10", "port":20000}}
  ]
}`
	str = `[{"path":"http://localhost/v1/aaa"}, {"path":"http://localhost:8080/v2/bbb", "delay":10}]`
	//j := make(map[string]interface{})
	var j interface{}
	if err := json.Unmarshal([]byte(str), &j); err != nil {
		fmt.Println("转换出错: ", err.Error())
		return
	}
	fmt.Println(j)
	fmt.Println("====================================")
	fmt.Println(str)
	fmt.Println("====================================")
	fmt.Println(j)
	//dump_map("", &j)
}

func (j *json_string) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("\"")
	buf.WriteString(*(*string)(j))
	buf.WriteString("\"")
	return buf.Bytes(), nil
}

func (j *json_string) UnmarshalJSON(d []byte) error {
	return nil
}

func (j *json_number) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("")
	switch v := j.value.(type) {
	case float32, float64, uint, uint32, uint64, int, int32, int64:
		buf.WriteString(fmt.Sprint(v))
	default:
		return nil, errors.New("not a number format")
	}
	return buf.Bytes(), nil
}

func (j *json_number) UnmarshalJSON(d []byte) error {
	return nil
}

func (o *Options) make_path(key ...string) *Options {
	ret := o
	return ret
}

func (o *Options) get_child_by_key(key string, create bool) *Options {
	return nil
}

func (o *Options) get_child_by_key2(create bool, path ...string) json_value {
	return nil
}

func (o *Options) get_value() (ret interface{}, exist bool) {
	return nil, false
}

func (o *Options) marshal_json(buf *bytes.Buffer) error {
	return nil
}

func (o *Options) parse_json_object(j *simplejson.Json) error {
	return nil
}

func (o *Options) parse_json_string(js string) error {
	return nil
}

func (o *Options) parse_ini_file() error {
	return nil
}

func (o *Options) parse_command_line() error {
	return nil
}
