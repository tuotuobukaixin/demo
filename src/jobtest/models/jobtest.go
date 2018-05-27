package models

import (
	"github.com/astaxie/beego/orm"
)

type Job struct {
	ID       int64  `orm:"pk;auto"`
	Name     string `orm:"size(64)"`
}


//AddVms Add a vms record
func AddJob(name string) error {
	var tmp = Job{}
	tmp.Name = name
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Insert(&tmp)
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return nil
}

func DeleteJob(id int) error {
	dao := orm.NewOrm()
	if err := dao.Begin(); err != nil {
		return err
	}
	_, err := dao.Raw("delete from job where  i_d = ?", id).Exec()
	if err == nil {
		dao.Commit()
	} else {
		dao.Rollback()
	}

	return err
}

func GetJob(name string) (*Job, error) {
	var engine Job
	dao := orm.NewOrm()

	err := dao.QueryTable("job").Filter("name", name).One(&engine)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func GetJobs() ([]Job, error) {
	var engine []Job
	dao := orm.NewOrm()

	_, err := dao.QueryTable("job").All(&engine)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
