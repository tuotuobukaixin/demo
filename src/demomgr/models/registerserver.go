package models

import (
	"github.com/astaxie/beego/orm"
)

type GameServer struct {
	ID          int64  `orm:"pk;auto"`
	Name        string `orm:"size(64)"`
	Status      string `orm:"size(64);null"`
	ServiceAddr string `orm:"size(64);null"`
	FileTest    bool   `orm:"default(false)"`
	FileSize    int    `orm:"size(32)"`
	TcpTest     bool   `orm:"default(false)"`
	TcpNum      int    `orm:"size(32)"`
}

type GameServerGet struct {
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	ServiceAddr string `json:"addr,omitempty"`
	FileTest    bool   `json:"filetest,omitempty"`
	FileSize    int    `json:"filesize,omitempty"`
	TcpTest     bool   `json:"tcptest,omitempty"`
	TcpNum      int    `json:"tcpnum,omitempty"`
}

//AddVms Add a vms record
func AddGameServer(end *GameServer) error {
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
func UpdateGameServer(end *GameServer) error {
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
	_, err := dao.Raw("delete from game_server where  i_d = ?", id).Exec()
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func GetGameServer(name string) (*GameServer, error) {
	var engine GameServer
	dao := orm.NewOrm()

	err := dao.QueryTable("game_server").Filter("name", name).One(&engine)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func GetGameServers() ([]GameServer, error) {
	var engine []GameServer
	dao := orm.NewOrm()

	_, err := dao.QueryTable("game_server").All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
