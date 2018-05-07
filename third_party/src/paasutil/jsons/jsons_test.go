package jsons

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func fload(filename string) (obj interface{}, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err == nil {
		err = load(&obj, data)
	}
	return
}
func TestDump(t *testing.T) {
	var obj interface{} = map[string]interface{}{
		"i": 10,
	}
	if data, err := DumpIndent(obj); err != nil {
		t.Fatal(err)
	} else if content := string(data); !strings.Contains(content, `"i": 10`) {
		t.Fatal(content)
	}
	objs := make([]interface{}, 0)
	objs = append(objs, obj)
	if data, err := Dump(objs); err != nil {
		t.Fatal(err)
	} else if content := string(data); !strings.HasPrefix(content, "[") || strings.HasSuffix(content, "] ") {
		t.Fatal(content)
	} else if err := Load(&obj, data); err != nil || len(obj.([]interface{})) != 1 {
		t.Fatal(obj, err)
	}
}

func TestGet(t *testing.T) {
	obj, _ := fload(".test/list.json")
	if price, ok := Get(obj, "tools.0/price"); !ok || price != float64(1) {
		t.Fatal(price)
	}
	if features, ok := Get(obj, "tools.0/features"); !ok || features.([]interface{})[0].(string) != "cut" {
		t.Fatal(features)
	} else {
		features.([]interface{})[0] = "run"
	}

	if features, ok := Get(obj, "tools.0/features"); !ok || features.([]interface{})[0].(string) != "run" {
		t.Fatal(features)
	}

	if self, found := Get(obj, ""); !found || self == nil {
		t.Fatal(self, found)
	} else if feature, found := GetString(self, "tools.0/features.0"); !found || feature != "run" {
		t.Fatal(feature, found)
	}

}

func TestGetMap(t *testing.T) {
	obj, _ := fload(".test/body.json")
	if id, found := GetString(obj, "access/user/id"); !found || id != "142801848c8c4ff38a98e24acfae901d" {
		t.Fatal(id)
	} else {
		if old, ok := Set(obj, "access/user/id", "hahaha"); !ok || old != "142801848c8c4ff38a98e24acfae901d" {
			t.Fatal(old)
		}
	}
	if user, found := GetMap(obj, "access/user"); !found {
		t.Fatal(found)
	} else {
		if id, found := GetString(user, "id"); !found || id != "hahaha" {
			t.Fatal(id)
		}
	}
}

func TestGetMap2(t *testing.T) {
	type mt map[string]interface{}
	me := mt{
		"me": "me",
	}
	obj := map[string]interface{}{}
	if old, ok := Set(obj, "me", me); !ok || old != nil {
		t.Fatal(obj, old, ok)
	}

	if me2, found := Get(obj, "me"); !found || me2 == nil || reflect.TypeOf(me2).Kind() != reflect.Map {
		t.Fatal(obj, me2, found)
	}

	if me2, found := GetMap(obj, "me"); !found {
		t.Fatal(obj, me2, found)
	} else if me2 == nil {
		t.Fatal(obj, me2, found)
	} else if val, found := GetString(me2, "me"); !found || val != "me" {
		t.Fatal(val, found)
	}

}

func TestGetSlice(t *testing.T) {
	obj, _ := fload(".test/body.json")
	services, _ := GetSlice(obj, "access/serviceCatalog")
	if len(services) == 0 {
		t.Fail()
	}
	name, _ := GetString(services[0], "name")
	if name != "nova" {
		t.Fatal(name)
	}
	if kind, _ := GetString(services[0], "type"); kind != "compute" {
		t.Fatal(kind)
	}
	idx := len(services)
	Put(obj, "access/serviceCatalog", "haha")
	haha, _ := GetString(obj, "access/serviceCatalog."+strconv.Itoa(idx))
	t.Log("access/serviceCatalog." + strconv.Itoa(idx))
	if haha != "haha" {
		t.Fatal(len(services))
	}
	if audit_ids, found := GetSlice(obj, "access/token/audit_ids"); !found {
		t.Fatal(found)
	} else if audit_ids[0].(string) != "lUUw0fPLSwO3jtUv22pSnQ" {
		t.Fatal(audit_ids)
	}
}

func TestGetSlice2(t *testing.T) {
	type msi map[string]interface{}
	obj := msi{}
	me := make([]msi, 0)
	me = append(me, msi{
		"a": "a",
	})
	obj["me"] = me
	if s, found := GetSlice(obj, "me"); !found || s == nil || len(s) != 1 {
		t.Fatal(s, found)
	} else if me2 := s[0].(map[string]interface{}); me2 == nil || me2["a"] != "a" {
		t.Fatal(me2)
	}
}

func TestGetRootSlice(t *testing.T) {
	obj := make([]interface{}, 0)
	obj = append(obj, "me")
	if slice, found := GetSlice(obj, ""); !found || len(slice) != 1 {
		t.Fatal(slice, found)
	}
	if old, ok := Put(obj, "", "me"); ok { // Do not support put to root slice
		t.Fatal(old, ok)
	}
	if me, found := GetString(obj, ".0"); !found || me != "me" {
		t.Fatal(me, found)
	}
	if old, ok := Set(obj, ".0", "you"); !ok || old != "me" {
		t.Fatal(old, ok)
	}
	if you, found := GetString(obj, ".0"); !found || you != "you" {
		t.Fatal(you, found)
	}
}

func TestGetBool(t *testing.T) {
	obj, _ := fload(".test/body.json")
	if enabled, found := GetBool(obj, "access/token/tenant/enabled"); !found || !enabled {
		t.Fatal(enabled)
	}
}

func TestGetNumber(t *testing.T) {
	obj, _ := fload(".test/body.json")
	if isAdmin, found := GetNumber(obj, "access/metadata/is_admin"); !found || isAdmin != float64(1) {
		t.Fatal(isAdmin)
	}
	if old, ok := Set(obj, "access/metadata/is_admin", 2); !ok || old != float64(1) {
		nval, _ := GetNumber(obj, "access/metadata/is_admin")
		t.Fatal(old, ok, reflect.TypeOf(old), nval)
	}
	if nval, found := GetNumber(obj, "access/metadata/is_admin"); !found || nval != float64(2) {
		t.Fatal(nval)
	}
}

func TestGetString(t *testing.T) {
	obj, _ := fload(".test/alarm.json")
	if id, found := GetString(obj, "alarm_id"); !found || id != "8b5ee570-03fe-44b9-982b-2f368df12f61" {
		t.Fatal(id)
	}
	if name, found := GetString(obj, "threshold_rule/meter_name"); !found || name != "instance" {
		t.Fatal(name)
	}

	if _, found := GetString(obj, "alarm_id/me"); found {
		t.Fatal(found)
	}

	if old, ok := Put(obj, "alarm_id/me", "me"); !ok || old != nil {
		t.Fatal(old, ok)
	}

	if me, found := GetString(obj, "alarm_id/me"); !found || me != "me" {
		t.Fatal(me, found)
	}
	if me, found := GetString(obj, "alarm_id/me.0"); found {
		t.Fatal(me, found)
	}
	if you, found := GetString(obj, "alarm_id/you.0"); found {
		t.Fatal(you, found)
	}
	if her, found := GetString(obj, "alarm_ids/her"); found {
		t.Fatal(her, found)
	}
}

func TestSetString(t *testing.T) {
	obj := make(map[string]interface{})
	Set(obj, "a/b/c", "me")
	if c, found := GetString(obj, "a/b/c"); !found || c != "me" {
		t.Fatal(c, found)
	}
	if b, found := GetMap(obj, "a/b"); !found || b == nil {
		t.Fatal(b, found)
	} else if c, found := GetString(b, "c"); !found || c != "me" {
		t.Fatal(c, found)
	}
	if old, ok := Set(obj, "a/b/c/d", make([]interface{}, 0)); !ok || old != nil {
		t.Fatal(old, ok)
	}
	if old, ok := Put(obj, "a/b/c/d", "d"); !ok || old != nil {
		t.Fatal(old, ok)
	}

	if old, ok := Set(obj, "a/b/c/d.0", "D"); !ok || old != "d" {
		t.Fatal(old, ok)
	}

	if old, ok := Set(obj, "a/b/c/d.1", "D2"); ok {
		t.Fatal(old, ok)
	}
}

func TestSetRoot(t *testing.T) {
	obj := make(map[string]interface{})
	if old, ok := Set(obj, "/", "me"); ok {
		t.Fatal(old, ok)
	}
	if old, ok := Set(obj, "/a/b/c", "me"); !ok || old != nil {
		t.Fatal(old, ok)
	}
}

func TestSetSlice(t *testing.T) {
	obj := make(map[string]interface{})
	Set(obj, "slice", make([]interface{}, 0))
	Put(obj, "slice", 1)
	Put(obj, "slice", 2)
	if slice, found := GetSlice(obj, "slice"); !found || len(slice) != 2 {
		t.Fatal(found)
	} else if s1, s2 := slice[0].(int), slice[1].(int); s1 != 1 || s2 != 2 {
		t.Fatal(slice)
	}
	if s, found := GetNumber(obj, "slice.0"); !found || s != 1 {
		t.Fatal(s, found)
	}

	if s, found := GetNumber(obj, "slice.1"); !found || s != 2 {
		t.Fatal(s, found)
	}

	if s, found := GetNumber(obj, "slice.2"); found {
		t.Fatal(s, found)
	}

	if s, found := GetNumber(obj, "slice.a"); found {
		t.Fatal(s, found)
	}

	if old, ok := Put(obj, "slice", make(map[string]interface{})); !ok || old != nil {
		t.Fatal(old, ok)
	}

	if me2, found := GetMap(obj, "slice.2"); !found || me2 == nil {
		t.Fatal(me2, found)
	}

	if old, ok := Set(obj, "slice.2/a/b/c", "me"); !ok || old != nil {
		t.Fatal(old, ok)
	}

	if c, found := GetString(obj, "slice.2/a/b/c"); !found || c != "me" {
		t.Fatal(c, found)
	}
}

func TestAlarmQuery(t *testing.T) {
	obj, _ := fload(".test/alarm.json")
	query, found := GetSlice(obj, "threshold_rule/query")
	if !found || len(query) != 2 {
		t.Log(query)
		t.Log(found)
		t.Fatal(query)
	}
	parseQuery := func() (result map[string]string) {
		result = make(map[string]string)
		for _, q := range query {
			field, _ := GetString(q, "field")
			value, _ := GetString(q, "value")
			result[field] = value
		}
		return result
	}
	result := parseQuery()
	if result["instance_id"] != "a-random_string_group-yv5p46ul6krj" {
		t.Fatal(result)
	}
	if result["stack_id"] != "f45b413d-081e-4a9e-a1ad-48e2afa0f601" {
		t.Fatal(result)
	}
}

func TestOthers(t *testing.T) {
	obj, _ := fload(".test/other.json")
	if bools, found := GetSlice(obj, "bools"); !found || len(bools) != 2 {
		t.Fatal(obj)
	} else if bools[0].(bool) != true || bools[1].(bool) != false {
		t.Fatal(bools)
	}
	if numbers, found := GetSlice(obj, "numbers"); !found || len(numbers) != 2 {
		t.Fatal()
	} else if numbers[0].(float64) != 1 || numbers[1].(float64) != 2 {
		t.Fatal()
	}
}
