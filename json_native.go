package json

import ()

type ()

func create_json_array(value ...interface{}) *[]interface{} {
	ret := make([]interface{}, 0, 10)
	for _, v := range value {
		ret = append(ret, v)
	}
	return &ret
}

func create_json_object(key string, value interface{}) *map[string]interface{} {
	ret := make(map[string]interface{})
	ret[key] = value
	return &ret
}

func (j *Json) get_value() (ret interface{}, exist bool) {
	ret = nil
	exist = false
	if d, ok := j.data.(*[]interface{}); ok {
		arr := []interface{}(*d)
		for _, v := range arr {
			switch use := v.(type) {
			case *[]interface{}:
			case *map[string]interface{}:
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
	case *map[string]interface{}:
		obj := get_object_child(d, key, false)
		if obj == nil {
			obj = new(map[string]interface{})
			(*d)[key] = obj
		}
		return &Json{obj}
	case *[]interface{}:
		obj := get_array_child(d, key, false)
		if obj == nil {
			obj = new(map[string]interface{})
			set_array(d, key, obj)
		}
		return &Json{obj}
	default:
		obj := new(map[string]interface{})
		j.data = create_json_array(d, create_json_object(key, obj))
		return &Json{obj}
	}
}

func (j *Json) get_child_by_key(key string) *Json {
	switch d := j.data.(type) {
	case *map[string]interface{}:
		obj := get_object_child(d, key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *[]interface{}:
				return &Json{obj}
			case *map[string]interface{}:
				return &Json{obj}
			default:
				return &Json{create_json_array(obj)}
			}
		}
	case *[]interface{}:
		obj := get_array_child(d, key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *[]interface{}:
				return &Json{obj}
			case *map[string]interface{}:
				return &Json{obj}
			default:
				return &Json{create_json_array(obj)}
			}
		}
	default:
		return nil
	}
}

func set_array(a *[]interface{}, key string, value interface{}) {
	d := ([]interface{})(*a)
	var use *map[string]interface{} = nil
	for _, u := range d {
		if m, ok := u.(*map[string]interface{}); ok {
			if use != nil {
				return
			}
			use = m
		}
	}
	if use != nil {
		(*use)[key] = value
	} else {
		append_array(a, create_json_object(key, value))
	}
}

func set_or_append_array(a *[]interface{}, key string, value interface{}) {
	d := ([]interface{})(*a)
	var use *map[string]interface{} = nil
	for _, u := range d {
		if m, ok := u.(*map[string]interface{}); ok {
			if use != nil {
				return
			}
			use = m
		}
	}
	if use != nil {
		append_object(use, key, value)
	} else {
		append_array(a, create_json_object(key, value))
	}
}

func append_array(a *[]interface{}, value ...interface{}) *[]interface{} {
	for _, v := range value {
		*a = append(*a, v)
	}
	return a
}

func append_object(o *map[string]interface{}, key string, value interface{}) {
	m := (*map[string]interface{})(o)
	if use, exist := (*m)[key]; exist {
		(*m)[key] = create_json_array(use, value)
	} else {
		(*m)[key] = value
	}
}

func get_object_child(src *map[string]interface{}, key string, create bool) interface{} {
	obj := (*map[string]interface{})(src)
	if data, exist := (*obj)[key]; exist {
		return data
	} else if create {
		ret := create_json_object(key, make(map[string]interface{}))
		(*obj)[key] = ret
		return ret
	}
	return nil
}

func get_array_child(src *[]interface{}, key string, create bool) interface{} {
	arr := ([]interface{})(*src)
	var use *map[string]interface{} = nil
	for _, c := range arr {
		if v, ok := c.(*map[string]interface{}); ok {
			if use != nil {
				return nil
			}
			use = v
		}
	}
	return get_object_child(use, key, create)
}

func transform_from_array(src []interface{}) *[]interface{} {
	ret := make([]interface{}, 0, 10)
	for _, v := range src {
		switch u := v.(type) {
		case []interface{}:
			ret = append(ret, transform_from_array(u))
		case map[string]interface{}:
			ret = append(ret, transform_from_map(u))
		default:
			ret = append(ret, v)
		}
	}
	return &ret
}

func transform_from_map(src map[string]interface{}) interface{} {
	ret := make(map[string]interface{})
	for k, v := range src {
		switch u := v.(type) {
		case []interface{}:
			ret[k] = transform_from_array(u)
		case map[string]interface{}:
			ret[k] = transform_from_map(u)
		default:
			ret[k] = v
		}
	}
	return &ret
}
