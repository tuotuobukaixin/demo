package jsons

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func dumpIndent(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "  ")
}

func dump(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func load(obj interface{}, data []byte) (err error) {
	err = json.Unmarshal(data, obj)
	return
}

func isKind(obj interface{}, kind reflect.Kind) bool {
	return reflect.TypeOf(obj).Kind() == kind
}

func isInt(obj interface{}) bool {
	return isKind(obj, reflect.Int)
}

func isFloat(obj interface{}) bool {
	return isKind(obj,
		reflect.Float64) || isKind(obj, reflect.Float32)
}

func isNumber(obj interface{}) bool {
	return isFloat(obj) || isInt(obj)
}

func isBool(obj interface{}) bool {
	return isKind(obj, reflect.Bool)
}

func isString(obj interface{}) bool {
	return isKind(obj, reflect.String)
}

func isMap(obj interface{}) bool {
	return isKind(obj, reflect.Map)
}

func isSlice(obj interface{}) bool {
	return isKind(obj, reflect.Slice)
}

func setMapIndex(obj interface{}, key string, value interface{}) {
	v := reflect.ValueOf
	v(obj).SetMapIndex(v(key), v(value))
}

func mapIndex(obj interface{}, key string) (value interface{}, found bool) {
	v := reflect.ValueOf
	val := v(obj).MapIndex(v(key))
	if val.IsValid() {
		value = val.Interface()
		found = true
	}
	return
}

// func mapKeys(obj interface{}) (keys []string) {
// 	v := reflect.ValueOf
// 	keys = make([]string, 0)
// 	for _, k := range v(obj).MapKeys() {
// 		keys = append(keys, k.String())
// 	}
// 	return
// }

func _get(obj interface{}, path string) (value interface{}, found bool) {
	if isSlice(obj) && strings.HasPrefix(path, ".") {
		objs := obj.([]interface{})
		if idx, err := strconv.Atoi(path[1:]); err == nil && idx < len(objs) {
			value = objs[idx]
			found = true
			return
		}
		return
	}
	if !isMap(obj) {
		return
	}
	start := strings.LastIndex(path, ".")
	if start == -1 {
		value, found = mapIndex(obj, path)
		return
	}
	id := path[0:start]
	if obj, ok := mapIndex(obj, id); ok {
		value, found = _get(obj, path[start:])
		return
	}
	return
}

func get(obj interface{}, path string) (value interface{}, found bool) {
	path = strings.Trim(path, "/")
	if path == "" {
		value = obj
		found = true
		return
	}
	var index int
	if index = strings.Index(path, "/"); index == -1 {
		value, found = _get(obj, path)
		return
	}
	id := path[0:index]
	path = path[index+1:]
	if obj, found = _get(obj, id); !found {
		return
	}
	value, found = get(obj, path)
	return
}

func _set(obj interface{}, path string, value interface{}) (old interface{}, ok bool) {
	if isSlice(obj) && strings.HasPrefix(path, ".") {
		objs := obj.([]interface{})
		if idx, err := strconv.Atoi(path[1:]); err == nil && idx < len(objs) {
			old = objs[idx]
			objs[idx] = value
			ok = true
			return
		}
		return
	}
	if !isMap(obj) {
		return
	}
	start := strings.LastIndex(path, ".")
	if start == -1 {
		if val, found := mapIndex(obj, path); found {
			old = val
		}
		setMapIndex(obj, path, value)
		ok = true
		return
	}
	id := path[0:start]
	if obj, found := mapIndex(obj, id); found {
		old, ok = _set(obj, path[start:], value)
		return
	}
	return
}

func set(obj interface{}, path string, value interface{}) (old interface{}, ok bool) {
	value = toValue(value)
	path = strings.Trim(path, "/")
	if path == "" {
		return
	}
	ids := strings.Split(path, "/")
	idx := len(ids) - 1
	if idx == 0 {
		old, ok = _set(obj, ids[0], value)
		return
	}

	for _, id := range ids[0:idx] {
		if m, found := GetMap(obj, id); found {
			obj = m
		} else {
			nm := map[string]interface{}{}
			if _, ok = _set(obj, id, nm); !ok {
				return
			}
			obj = nm
		}
	}
	old, ok = _set(obj, ids[idx], value)
	return
}

func put(obj interface{}, path string, value interface{}) (old interface{}, ok bool) {
	if slice, found := GetSlice(obj, path); found {
		value = append(slice, value)
		if path == "" {
			ok = false
			return
		}
		_, ok = set(obj, path, value)
		return
	}
	old, ok = set(obj, path, value)
	return
}

var (
	Get        = get
	Set        = set
	Put        = put
	Dump       = dump
	DumpIndent = dumpIndent
	Load       = load
)

func toBool(obj interface{}) (value bool) {
	switch obj.(type) {
	case bool:
		value = obj.(bool)
	}
	return
}

func toString(obj interface{}) (value string) {
	switch obj.(type) {
	case string:
		value = obj.(string)
	}
	return
}

func toNumber(obj interface{}) (value float64) {
	v := reflect.ValueOf(obj)
	switch obj.(type) {
	case int, int8, int16, int32, int64:
		value = float64(v.Int())
	case float32, float64:
		value = v.Float()
	}
	return
}

func toMap(obj interface{}) (value map[string]interface{}) {
	switch obj.(type) {
	case map[string]interface{}:
		value = obj.(map[string]interface{})
	default:
		v := reflect.ValueOf
		t := reflect.TypeOf
		value = v(obj).Convert(t(value)).Interface().(map[string]interface{})
	}
	return
}

func toSlice(obj interface{}) (value []interface{}) {
	switch obj.(type) {
	case []interface{}:
		value = obj.([]interface{})
	default:
		if value = make([]interface{}, 0); obj != nil {
			v := reflect.ValueOf
			val := v(obj)
			for i := 0; i < val.Len(); i++ {
				value = append(value, toValue(val.Index(i).Interface()))
			}
		}
	}
	return
}

func toValue(obj interface{}) (value interface{}) {
	if isMap(obj) {
		value = toMap(obj)
	} else if isSlice(obj) {
		value = toSlice(obj)
	} else {
		value = obj
	}
	return
}
func GetBool(obj interface{}, path string) (value bool, found bool) {
	if i, ok := get(obj, path); ok && isBool(i) {
		value = toBool(i)
		found = true
	}
	return
}
func GetNumber(obj interface{}, path string) (value float64, found bool) {
	if i, ok := get(obj, path); ok && isNumber(i) {
		value = toNumber(i)
		found = true
	}
	return
}

func GetString(obj interface{}, path string) (value string, found bool) {
	if i, ok := get(obj, path); ok && isString(i) {
		value = toString(i)
		found = true
	}
	return
}

func GetSlice(obj interface{}, path string) (value []interface{}, found bool) {
	if i, ok := get(obj, path); ok && isSlice(i) {
		value = toSlice(i)
		found = true
	}
	return
}

func GetMap(obj interface{}, path string) (value map[string]interface{}, found bool) {
	if i, ok := get(obj, path); ok && isMap(i) {
		value = toMap(i)
		found = true
	}
	return
}
