package json

import ()

type (
	json_value  interface{}
	json_string string
	json_number struct {
		value interface{}
	}
	json_bool bool

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
	ret = nil
	exist = false
	if d, ok := j.data.(*json_array); ok {
		arr := []json_value(*d)
		for _, v := range arr {
			switch use := v.(type) {
			case *json_array:
			case *json_object:
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
	case *json_object:
		obj := d.get_child_by_key(key, false)
		if obj == nil {
			obj = &json_object{}
			d.set(key, obj)
		}
		return &Json{obj}
	case *json_array:
		obj := d.get_child_by_key(key, false)
		if obj == nil {
			obj = &json_object{}
			d.set(key, obj)
		}
		return &Json{obj}
	default:
		obj := &json_object{}
		j.data = create_json_array(d, create_json_object(key, obj))
		return &Json{obj}
	}
}

func (j *Json) get_child_by_key(key string) *Json {
	switch d := j.data.(type) {
	case *json_object:
		obj := d.get_child_by_key(key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *json_array:
				return &Json{obj}
			case *json_object:
				return &Json{obj}
			default:
				return &Json{create_json_array(obj)}
			}
		}
	case *json_array:
		obj := d.get_child_by_key(key, false)
		if obj == nil {
			return nil
		} else {
			switch obj.(type) {
			case *json_array:
				return &Json{obj}
			case *json_object:
				return &Json{obj}
			default:
				return &Json{create_json_array(obj)}
			}
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
	var use *json_object = nil
	for _, c := range arr {
		if v, ok := c.(*json_object); ok {
			if use != nil {
				return nil
			}
			use = v
		}
	}
	return use.get_child_by_key(key, create)
}

func transform_from_array(src []interface{}) *json_array {
	ret := &json_array{}
	for _, v := range src {
		switch u := v.(type) {
		case []interface{}:
			ret.append(transform_from_array(u))
		case map[string]interface{}:
			ret.append(transform_from_map(u))
		default:
			ret.append(v)
		}
	}
	return ret
}

func transform_from_map(src map[string]interface{}) json_value {
	ret := &json_object{}
	for k, v := range src {
		switch u := v.(type) {
		case []interface{}:
			ret.set(k, transform_from_array(u))
		case map[string]interface{}:
			ret.set(k, transform_from_map(u))
		default:
			ret.set(k, v)
		}
	}
	return ret
}
