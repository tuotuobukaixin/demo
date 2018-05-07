package jsons

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSet(t *testing.T) {
	obj := NewJsonObj()
	obj.Set("a", "v")
	obj.Get("a")
	if a := obj.GetString("a"); a != "v" {
		t.Fatal(a)
	}
}

func TestPut(t *testing.T) {
	obj := NewJsonObj()
	obj.Set("a", make([]interface{}, 0))
	obj.Put("a", "v")
	if a := obj.GetSlice("a"); len(a) == 0 || a[0].(string) != "v" {
		t.Fatal(a)
	}
}

func TestLoadJsonObj(t *testing.T) {
	obj := LoadJsonObj(".test/alarm.json")
	if action := obj.GetString("alarm_actions.0"); !strings.HasPrefix(action, "http://158.85.90.245:8000") {
		t.Fatal(action)
	}
	if enabled := obj.GetBool("enabled"); !enabled {
		t.Fatal(enabled)
	}
	if period := obj.GetNumber("threshold_rule/period"); period != 60 {
		t.Fatal(period)
	}
	if query := obj.GetSlice("threshold_rule/query"); len(query) != 2 {
		t.Fatal(query)
	}
	if rule := obj.GetJsonObj("threshold_rule"); rule == nil {
		t.Fatal(rule)
	} else if mname := rule.GetString("meter_name"); mname != "instance" {
		t.Fatal(mname)
	}
}

func TestSaveJsonObj(t *testing.T) {
	obj := NewJsonObj()
	filename := filepath.Join(os.TempDir(), "empty.json")
	defer os.Remove(filename)
	if err := SaveJsonObj(filename, obj); err != nil {
		t.Fatal(err)
	}
	if data, err := ioutil.ReadFile(filename); err != nil || string(data) != "{}" {
		t.Fatal(data, err, string(data))
	}
	obj.Set("a/b/c", "me")
	if err := SaveJsonObj(filename, obj); err != nil {
		t.Fatal(err)
	}
	if data, err := ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	} else if content := string(data); !strings.Contains(content, `"c":"me"`) {
		t.Fatal(content)
	} else if content != obj.Dump() {
		t.Fatal(content, obj.Dump())
	}
}

func TestHas(t *testing.T) {
	obj := NewJsonObj()
	obj.Set("a/b/c", "me")
	if has := obj.Has("a/b/c"); !has {
		t.Fatal(has)
	}
	if has := obj.Has("/a/b"); !has {
		t.Fatal(has)
	}
	if has := obj.Has("/a"); !has {
		t.Fatal(has)
	}
	if has := obj.Has("/"); !has {
		t.Fatal(has)
	}
}

func TestSetJsonObj(t *testing.T){
	obj := NewJsonObj()
	me := NewJsonObj()
	me.Set("a/b/c", "me")
	obj.Set("me", me)
	if s := obj.GetJsonObj("me"); s == nil {
		t.Fatal(s, me)
	}
}

func TestAlarmActions(t *testing.T) {
	alarm := LoadJsonObj(".test/alarm.json")
	if alarm_action := alarm.GetString("alarm_actions.0"); alarm_action == "" {
		t.Fatal(alarm_action)
	}
}

func TestNewJsonObj(t *testing.T) {
	if obj := NewJsonObj(); len(obj) != 0 {
		t.Fatal(obj)
	}
	content1 := `{"a":"b"}`
	if obj := NewJsonObj(content1); obj.IsEmpty(){
		t.Fatal(obj)
	}else if obj.GetString("a") != "b" {
		t.Fatal(obj)
	}
	content2 := `{"c":"d"}`
	if obj := NewJsonObj(content1, content2); obj.IsEmpty(){
		t.Fatal(obj)
	}else if obj.GetString("c") != "d" {
		t.Fatal(obj)
	}
}

func TestIsEmpty(t *testing.T) {
	var obj1 JsonObj
	if empty := obj1.IsEmpty(); !empty {
		t.Fatal(empty)
	}
	obj2 := NewJsonObj()
	if empty := obj2.IsEmpty(); !empty {
		t.Fatal(empty)
	}
	obj3 := JsonObj{}
	if empty := obj3.IsEmpty(); !empty {
		t.Fatal(empty)
	}
	
	obj3.Set("a", "av")
	if empty := obj3.IsEmpty(); empty {
		t.Fatal(empty)
	}
	var obj4 JsonObj = nil
	if empty := obj4.IsEmpty(); !empty {
		t.Fatal(empty)
	}	
}

func TestJsonObjLoad(t *testing.T){
	obj := NewJsonObj()
	if err := obj.Load("") ; err != nil {
		t.Fatal(obj)
	}else if !obj.IsEmpty() {
		t.Fatal(obj)
	}
	obj = LoadJsonObj(".test/body.json")
	obj2 := JsonObj{}
	obj2.Load(obj.Dump())
	if obj2.IsEmpty(){
		t.Log(obj)
		t.Fatal(obj2)
	}
	obj3 := jsonobj_test_jsonobj()
	obj3.Load(obj.Dump())
	if !obj3.IsEmpty(){
		t.Log(obj)
		t.Fatal(obj3)
	}
}

func jsonobj_test_jsonobj() (obj JsonObj){
	return
}
