package models

import (
	"github.com/astaxie/beego/orm"
)

type DemoTest struct {
	ID          int64  `orm:"pk;auto"`
	Name        string `orm:"size(64)"`
	Status      string `orm:"size(64);null"`
	ServiceAddr string `orm:"size(64);null"`
	Podip       string `orm:"size(64);null"`
	FileTest    bool   `orm:"default(false)"`
	FileSize    int    `orm:"size(32)"`
	TcpTest     bool   `orm:"default(false)"`
	TcpNum      int    `orm:"size(32)"`
}

type DemoTestGet struct {
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	ServiceAddr string `json:"addr,omitempty"`
	Podip   string `json:"podip,omitempty"`
	FileTest    bool   `json:"filetest,omitempty"`
	FileSize    int    `json:"filesize,omitempty"`
	TcpTest     bool   `json:"tcptest,omitempty"`
	TcpNum      int    `json:"tcpnum,omitempty"`
}

//AddVms Add a vms record
func AddDemoTest(end *DemoTest) error {
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
func UpdateDemoTest(end *DemoTest) error {
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

func DeleteDemoTest(id int) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Raw("delete from demo_test where  i_d = ?", id).Exec()
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func GetDemoTest(name string) (*DemoTest, error) {
	var engine DemoTest
	dao := orm.NewOrm()

	err := dao.QueryTable("demo_test").Filter("name", name).One(&engine)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func GetDemoTests() ([]DemoTest, error) {
	var engine []DemoTest
	dao := orm.NewOrm()

	_, err := dao.QueryTable("demo_test").All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
