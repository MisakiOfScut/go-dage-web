package dao

import (
	"fmt"
	"go-dage-web/app/utils/mysql"
	"go.uber.org/zap"
	"time"
)

func init(){
	timeLocal := time.FixedZone("CST", 3600*8)
	time.Local = timeLocal
}

const NodeListTable = "node_list"

type ListItem struct {
	Name           string `db:"script_name" json:"script_name"`
	Date           string `db:"date" json:"date"`
	PublishVersion int    `db:"publish_version" json:"publish_version"`
}
type List []ListItem

func (p *List) Get() error {
	sqlStmt := "SELECT * from " + NodeListTable
	if err := mysql.GSqlDB.Select(p, sqlStmt); err != nil {
		zap.S().Errorf("sql:%s, err:%v", sqlStmt, err)
		return err
	}

	zap.S().Infof("sql:%s, query %d nodes from db", sqlStmt, len(*p))
	return nil
}

func (p *ListItem) Add() (bool, error) {
	sqlStmt := fmt.Sprintf("INSERT INTO %s (`script_name`, `date`) ", NodeListTable) +
		fmt.Sprintf("SELECT * FROM (SELECT '%s' AS script_name, ? AS date) AS tmp ", p.Name) +
		fmt.Sprintf("WHERE NOT EXISTS (SELECT script_name FROM %s WHERE script_name = '%s') LIMIT 1;", NodeListTable,
			p.Name)

	if res, err := mysql.GSqlDB.Exec(sqlStmt, time.Now().Local()); err != nil {
		zap.S().Errorf("sql:%s, err: %v", sqlStmt, err)
		return false, err
	} else if affected, err := res.RowsAffected(); err != nil {
		zap.S().Errorf("sql:%s, res.RowsAffected() return err: %v", sqlStmt, err)
		return false, err
	} else {
		zap.S().Debugf("insert script_name:%s and %d row affected", p.Name, affected)
		return affected == 1, nil
	}
}

func (p *ListItem) Delete() error {
	sqlStmt := fmt.Sprintf("DELETE From %s where script_name=?", NodeListTable)
	result, err := mysql.GSqlDB.Exec(sqlStmt, p.Name)
	if err != nil {
		zap.S().Errorf("sql:%s, err: %v", sqlStmt, err)
		return fmt.Errorf("删除配置%s失败", p.Name)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("get affected failed, err:%v\n", err)
		return fmt.Errorf("删除配置%s失败", p.Name)
	}
	zap.S().Debugf("delete script_name:%s and %d row affected", p.Name, affectedRows)
	return nil
}
