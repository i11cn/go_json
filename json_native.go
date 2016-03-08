package json

import ()

type ()

func create_json_array(value ...interface{}) *[]interface{} {
	ret := make([]interface{}, 0, 10)
	ret = append(ret, value...)
	return &ret
}

func create_json_object(key string, value interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	ret[key] = value
	return ret
}

func (j *Json) get_value() (ret interface{}, exist bool) {
	ret = nil
	exist = false
	if d, ok := j.data.(*[]interface{}); ok {
		for _, v := range *d {
			switch use := v.(type) {
			case *[]interface{}:
			case map[string]interface{}:
			case []interface{}:
				j.data = &j.data
			default:
				ret = use
				exist = true
				return
			}
		}
	}
	return
}

func (j *Json) get_or_create_child(key string) *Json {
	switch d := j.data.(type) {
	case map[string]interface{}:
		obj := get_object_child(d, key, false)
		if obj == nil {
			obj = make(map[string]interface{})
			d[key] = obj
		}
		return &Json{obj}
	case *[]interface{}:
		obj := get_array_child(*d, key, false)
		if obj == nil {
			obj = create_json_object(key, make(map[string]interface{}))
			*d = append(*d, obj)
		}
		return &Json{obj}
	case []interface{}:
		j.data = &j.data
		d2, _ := j.data.(*[]interface{})
		obj := get_array_child(*d2, key, false)
		if obj == nil {
			obj = create_json_object(key, make(map[string]interface{}))
			*d2 = append(*d2, obj)
		}
		return &Json{obj}
	default:
		obj := make(map[string]interface{})
		j.data = create_json_array(d, create_json_object(key, obj))
		return &Json{obj}
	}
}

func (j *Json) get_child_by_key(key string) *Json {
	switch d := j.data.(type) {
	case map[string]interface{}:
		obj := get_object_child(d, key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *[]interface{}:
				return &Json{obj}
			case map[string]interface{}:
				return &Json{obj}
			default:
				return &Json{create_json_array(obj)}
			}
		}
	case *[]interface{}:
		obj := get_array_child(*d, key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *[]interface{}:
				return &Json{obj}
			case map[string]interface{}:
				return &Json{obj}
			default:
				return &Json{create_json_array(obj)}
			}
		}
	case []interface{}:
		j.data = &j.data
		obj := get_array_child(d, key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *[]interface{}:
				return &Json{obj}
			case map[string]interface{}:
				return &Json{obj}
			case []interface{}:
				return &Json{&obj}
			default:
				return &Json{create_json_array(obj)}
			}
		}
	default:
		return nil
	}
}

func set_array(a []interface{}, key string, value interface{}) []interface{} {
	var use map[string]interface{} = nil
	for _, u := range a {
		if m, ok := u.(map[string]interface{}); ok {
			if use != nil {
				return a
			}
			use = m
		}
	}
	if use != nil {
		use[key] = value
	} else {
		a = append(a, create_json_object(key, value))
	}
	return a
}

func set_or_append_array(a []interface{}, key string, value interface{}) []interface{} {
	var use map[string]interface{} = nil
	for _, u := range a {
		if m, ok := u.(map[string]interface{}); ok {
			if use != nil {
				return a
			}
			use = m
		}
	}
	if use != nil {
		use[key] = value
	} else {
		a = append(a, create_json_object(key, value))
	}
	return a
}

func append_object(o map[string]interface{}, key string, value interface{}) map[string]interface{} {
	if use, exist := o[key]; exist {
		o[key] = create_json_array(use, value)
	} else {
		o[key] = value
	}
	return o
}

func get_object_child(src map[string]interface{}, key string, create bool) interface{} {
	if data, exist := src[key]; exist {
		return data
	} else if create {
		ret := create_json_object(key, make(map[string]interface{}))
		src[key] = ret
		return ret
	}
	return nil
}

func get_array_child(src []interface{}, key string, create bool) interface{} {
	var use map[string]interface{} = nil
	for _, c := range src {
		if v, ok := c.(map[string]interface{}); ok {
			if use != nil {
				return nil
			}
			use = v
		}
	}
	return get_object_child(use, key, create)
}
