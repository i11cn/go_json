package json

import (
	"encoding/json"
	"fmt"
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
)

func create_json_array(value ...interface{}) *json_array {
	ret := json_array{}
	for _, v := range value {
		ret = append(ret, v)
	}
	return &ret
}

func create_json_object(key string, value interface{}) *json_object {
	ret := map[string]json_value{}
	ret[key] = value
	return (*json_object)(&ret)
}

func (j *Json) get_value() (ret interface{}, exist bool) {
	ret = j.data
	exist = j.data != nil
	return
}

func (j *Json) get_or_create_child(key string) *Json {
	if j == nil {
		return NewJson()
	}
	switch d := j.data.(type) {
	case nil:
		use := create_json_object(key, &json_object{})
		j.data = use
		return j.get_child_by_key(key)
	case *json_object:
		return &Json{d.get_child_by_key(key, true)}
	case *json_array:
		return &Json{d.get_child_by_key(key, true)}
	default:
		return nil
	}
	return NewJson()
}

func (j *Json) get_child_by_key(key string) *Json {
	if j == nil {
		return nil
	}
	switch d := j.data.(type) {
	case nil:
		return nil
	case *json_object:
		obj := d.get_child_by_key(key, false)
		if obj == nil {
			return nil
		} else {
			return &Json{obj}
		}
	case *json_array:
		obj := d.get_child_by_key(key, false)
		if obj == nil {
			return nil
		} else {
			return &Json{}
		}
	default:
		return nil
	}
}

func (a *json_array) set(key string, value interface{}) {
	d := ([]json_value)(*a)
	var use *json_object = nil
	for _, u := range d {
		if m, ok := u.(*json_object); ok {
			if use != nil {
				return
			}
			use = m
		}
	}
	if use != nil {
		use.set(key, value)
	} else {
		a.append(create_json_object(key, value))
	}
}

func (a *json_array) set_or_append(key string, value interface{}) {
	d := ([]json_value)(*a)
	var use *json_object = nil
	for _, u := range d {
		if m, ok := u.(*json_object); ok {
			if use != nil {
				return
			}
			use = m
		}
	}
	if use != nil {
		use.append(key, value)
	} else {
		a.append(create_json_object(key, value))
	}
}

func (a *json_array) append(value ...interface{}) *json_array {
	for _, v := range value {
		*a = append(*a, v)
	}
	return a
}

func (o *json_object) set(key string, value interface{}) {
	use := (*map[string]json_value)(o)
	(*use)[key] = value
}

func (o *json_object) append(key string, value interface{}) {
	m := (*map[string]json_value)(o)
	if use, exist := (*m)[key]; exist {
		(*m)[key] = create_json_array(use, value)
	} else {
		(*m)[key] = value
	}
}

func (src *json_object) get_child_by_key(key string, create bool) json_value {
	obj := (*map[string]json_value)(src)
	if data, exist := (*obj)[key]; exist {
		return data
	} else if create {
		ret := create_json_object(key, json_object{})
		(*obj)[key] = ret
		return ret
	}
	return nil
}

func (src *json_array) get_child_by_key(key string, create bool) json_value {
	arr := ([]json_value)(*src)
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
