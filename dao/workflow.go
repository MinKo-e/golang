package dao

import (
	"errors"
	"github.com/sirupsen/logrus"
	db "k8sManager/db/mysql"
	"k8sManager/model"
)

var Workflow workflow

type workflow struct {
}

type WorkResp struct {
	Items []*model.Workflow `json:"items"`
	Total int               `json:"total"`
}

func (w *workflow) GetList(name string, limit, page int) (*WorkResp, error) {
	startIndex := (page - 1) * limit
	var workflowList []*model.Workflow
	tx := db.DB.Where("name like ?", "%"+name+"%").Limit(limit).Offset(startIndex).Order("id desc").Find(&workflowList)

	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logrus.Error("获取workflow列表失败" + tx.Error.Error())
		return nil, errors.New("获取workflow列表失败" + tx.Error.Error())
	}
	return &WorkResp{
		Items: workflowList,
		Total: len(workflowList),
	}, nil
}

func (w *workflow) GetByID(id int) (*model.Workflow, error) {
	var data model.Workflow
	tx := db.DB.Where("id = ?", id).First(&data)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logrus.Error("获取workflow列表失败" + tx.Error.Error())
		return &model.Workflow{}, errors.New("获取workflow列表失败" + tx.Error.Error())
	}
	logrus.Info("查询数据成功!")
	return &data, nil
}

func (w *workflow) DeleteId(id int) (err error) {
	tx := db.DB.Where("id = ?", id).Delete(model.Workflow{})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logrus.Error("删除workflow失败" + err.Error())
		return errors.New("删除workflow失败" + err.Error())
	}
	logrus.Info("删除数据成功!")
	return nil
}
func (w *workflow) Create(WS *model.Workflow) (err error) {
	tx := db.DB.Create(WS)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logrus.Error("创建workflow失败" + err.Error())
		return errors.New("创建workflow失败" + err.Error())
	}
	logrus.Info("创建数据成功")
	return nil
}
