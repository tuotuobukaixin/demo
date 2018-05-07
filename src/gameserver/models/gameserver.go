package models

import (
	"github.com/astaxie/beego/orm"
)

type GameServerTestResult struct {
	ID             int64  `orm:"pk;auto"`
	Name           string `orm:"size(64)"`
	Time           string `orm:"size(64)"`
	FileWriteSpeed int    `orm:"size(32)"`
	FileReadSpeed  int    `orm:"size(32)"`
	Success        int    `orm:"size(32)"`
	Total          int    `orm:"size(32)"`
}
type GameServerTestResultGet struct {
	Name           string `json:"name,omitempty"`
	Time           string `json:"time,omitempty"`
	FileWriteSpeed int    `json:"writespeed,omitempty"`
	FileReadSpeed  int    `json:"readspeed,omitempty"`
	Success        int    `json:"success,omitempty"`
	Total          int    `json:"total,omitempty"`
}

//AddVms Add a vms record
func AddGameServerTestResult(end *GameServerTestResult) error {
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
func UpdateGameServer(end *GameServerTestResult) error {
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

func DeleteGameServer(id int) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Raw("delete from game_server_test_result where  i_d = ?", id).Exec()
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func GetGameServer(name string) (*GameServerTestResult, error) {
	var engine GameServerTestResult
	dao := orm.NewOrm()

	err := dao.QueryTable("game_server_test_result").Filter("name", name).One(&engine)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func GetGameServers() ([]GameServerTestResult, error) {
	var engine []GameServerTestResult
	dao := orm.NewOrm()

	_, err := dao.QueryTable("game_server_test_result").All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
