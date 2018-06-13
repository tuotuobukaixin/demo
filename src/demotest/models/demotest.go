package models

import (
	"github.com/astaxie/beego/orm"
)

type DemoTestTestResult struct {
	ID             int64  `orm:"pk;auto"`
	Name           string `orm:"size(64)"`
	Time           string `orm:"size(64)"`
	FileWriteSpeed int    `orm:"size(32)"`
	FileReadSpeed  int    `orm:"size(32)"`
	Success        int    `orm:"size(32)"`
	Total          int    `orm:"size(32)"`
	Detail         string `orm:"type(text);null"`
}
type DemoTestTestResultGet struct {
	Name           string `json:"name,omitempty"`
	Time           string `json:"time,omitempty"`
	FileWriteSpeed int    `json:"writespeed,omitempty"`
	FileReadSpeed  int    `json:"readspeed,omitempty"`
	Success        int    `json:"success,omitempty"`
	Total          int    `json:"total,omitempty"`
}

//AddVms Add a vms record
func AddDemoTestTestResult(end *DemoTestTestResult) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Insert(end)
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return nil
}

//UpdateRuntimeEngine Update K8sRuntime record
func UpdateDemoTest(end *DemoTestTestResult) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Update(end)
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func DeleteDemoTest(name string) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Raw("delete from demo_test_test_result where  name = ?", name).Exec()
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func GetDemoTestResult(name string) ([]DemoTestTestResult, error) {
	var engine []DemoTestTestResult
	dao := orm.NewOrm()

	_,err := dao.QueryTable("demo_test_test_result").Filter("name", name).All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func GetDemoTestsResult() ([]DemoTestTestResult, error) {
	var engine []DemoTestTestResult
	dao := orm.NewOrm()

	_, err := dao.QueryTable("demo_test_test_result").All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
