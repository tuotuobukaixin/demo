package jsons

import (
	"io/ioutil"
)

type JsonObj map[string]interface{}

func NewJsonObj(contents ...string) (value JsonObj) {
	value = make(JsonObj)
	if len(contents) == 0 {
		return
	}
	for _, content := range contents {
		value.Load(content)
	}
	return
}

func GetJsonObj(json interface{}) (obj JsonObj) {
	switch json.(type) {
	case string:
		body := json.(string)
		if body != "" {
			obj = NewJsonObj(body)
		}
	case map[string]interface{}:
		obj = json.(map[string]interface{})
	case JsonObj:
		obj = json.(JsonObj)
	}
	return
}

func LoadJsonObj(filename string) (obj JsonObj) {
	if data, err := ioutil.ReadFile(filename); err == nil {
		obj = make(JsonObj)
		obj.Load(string(data))
	}
	return
}

func SaveJsonObj(filename string, obj JsonObj) (err error) {
	var data []byte
	if data, err = Dump(obj); err == nil {
		err = ioutil.WriteFile(filename, data, 0600)
	}
	return
}

func (self JsonObj) Dump() (str string) {
	if data, err := Dump(self); err == nil {
		str = string(data)
	}
	return
}

func (self JsonObj) String() (str string) {
	if data, err := DumpIndent(self); err == nil {
		str = string(data)
	}
	return
}

func (self JsonObj) Load(content string) (err error) {
	if content == "" {
		return
	}
	data := []byte(content)
	err = Load(&self, data)
	return
}

func (self JsonObj) Get(path string) (value interface{}) {
	if v, found := Get(self, path); found {
		value = v
	}
	return
}

func (self JsonObj) Set(path string, value interface{}) (old interface{}, ok bool) {
	old, ok = Set(self, path, value)
	return
}

func (self JsonObj) Put(path string, value interface{}) (old interface{}, ok bool) {
	old, ok = Put(self, path, value)
	return
}

func (self JsonObj) GetString(path string) (value string) {
	if v, ok := GetString(self, path); ok {
		value = v
	}
	return
}

func (self JsonObj) GetBool(path string) (value bool) {
	if v, ok := GetBool(self, path); ok {
		value = v
	}
	return
}

func (self JsonObj) GetNumber(path string) (value float64) {
	if v, ok := GetNumber(self, path); ok {
		value = v
	}
	return
}

func (self JsonObj) GetJsonObj(path string) (value JsonObj) {
	if v, found := GetMap(self, path); found {
		value = v
	}
	return
}

func (self JsonObj) GetSlice(path string) (value []interface{}) {
	if v, found := GetSlice(self, path); found {
		value = v
	}
	return
}

func (self JsonObj) Has(path string) (ok bool) {
	if _, found := Get(self, path); found {
		ok = true
	}
	return
}

func (self JsonObj) IsEmpty() (empty bool) {
	empty = (self == nil || len(self) == 0)
	return
}
