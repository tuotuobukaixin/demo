package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	ID       int64  `orm:"pk;auto"`
	Name     string `orm:"size(64)"`
	Password string `orm:"size(64);null"`
	Role     string `orm:"size(64);null"`
}

type UserPost struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

//AddVms Add a vms record
func AddUser(end *User) error {
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
func UpdateUser(end *User) error {
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

func DeleteUser(id int) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Raw("delete from user where  i_d = ?", id).Exec()
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func GetUser(name string) (*User, error) {
	var engine User
	dao := orm.NewOrm()

	err := dao.QueryTable("user").Filter("name", name).One(&engine)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func GetGameUser() ([]User, error) {
	var engine []User
	dao := orm.NewOrm()

	_, err := dao.QueryTable("user").All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
