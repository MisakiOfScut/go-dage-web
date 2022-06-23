package dao

import (
	"fmt"
	"go-dage-web/app/utils/mysql"
	"go.uber.org/zap"
	"time"
)

const (
	ScriptListTable = "script_list"
)

type Script struct {
	Name    string `db:"script_name" json:"script_name"`
	Content string `db:"content" json:"content"`
	Version int    `db:"version" json:"version"`
	Date    string `db:"date" json:"date"`
}

type Scripts []Script

func (p *Script) Get() error {
	sqlStmt := fmt.Sprintf("SELECT content,date FROM %s WHERE script_name=? AND version=?;", ScriptListTable)
	if err := mysql.GSqlDB.Get(p, sqlStmt, p.Name, p.Version); err != nil {
		zap.S().Errorf("sql:%s, err:%v", sqlStmt, err)
		return err
	}

	zap.S().Infof("sql:%s, query script %s, version %d from db", sqlStmt, p.Name, p.Version)
	return nil
}

func (p *Script) Add() error {
	scripts := Scripts{}
	sqlStmt := fmt.Sprintf("SELECT version FROM %s WHERE script_name=?;", ScriptListTable)
	if err := mysql.GSqlDB.Select(&scripts, sqlStmt, p.Name); err != nil {
		zap.S().Errorf("sql:%s, err:%v", sqlStmt, err)
		return fmt.Errorf("query script %s failed", p.Name)
	}
	// find the biggest version
	biggest := 0
	for i, _ := range scripts {
		if scripts[i].Version > biggest {
			biggest = scripts[i].Version
		}
	}
	if p.Version <= biggest {
		return fmt.Errorf("script:%s, 用户端当前最大版本%d已不是最新版本，请刷新页面", p.Name, p.Version-1)
	}

	sqlStmt = fmt.Sprintf("INSERT INTO %s ", ScriptListTable) +
		fmt.Sprintf("SELECT * FROM (SELECT '%s' AS script_name, ? AS content, %d AS version, ? AS date) AS tmp ",
			p.Name, p.Version) +
		fmt.Sprintf("WHERE NOT EXISTS (SELECT script_name FROM %s WHERE script_name='%s' AND version=%d) LIMIT 1;",
			ScriptListTable,
			p.Name, p.Version)

	if res, err := mysql.GSqlDB.Exec(sqlStmt, p.Content, time.Now()); err != nil {
		zap.S().Errorf("sql:%s, err: %v", sqlStmt, err)
		return fmt.Errorf("保存失败")
	} else if affected, err := res.RowsAffected(); err != nil {
		zap.S().Errorf("sql:%s, res.RowsAffected() return err: %v", sqlStmt, err)
		return fmt.Errorf("保存失败")
	} else {
		zap.S().Debugf("insert script_name:%s, version:%d and %d row affected", p.Name, p.Version, affected)
		if affected == 0 {
			return fmt.Errorf("script:%s, 用户端当前最大版本%d已不是最新版本，请刷新页面", p.Name, p.Version-1)
		}
		return nil
	}
}

func (p *Scripts) GetAll(name string) error {
	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE script_name=? ORDER BY version;", ScriptListTable)
	if err := mysql.GSqlDB.Select(p, sqlStmt, name); err != nil {
		zap.S().Errorf("sql:%s, err:%v", sqlStmt, err)
		return err
	}

	zap.S().Infof("sql:%s, query script %s with %d versions from db", sqlStmt, name, len(*p))
	return nil
}

func (p *Script) Publish() error {
	sqlStmt := fmt.Sprintf("UPDATE %s SET publish_version=?, date=? WHERE script_name=?;", NodeListTable)
	if res, err := mysql.GSqlDB.Exec(sqlStmt, p.Version, time.Now(), p.Name); err != nil {
		zap.S().Errorf("sql:%s, err:%v", sqlStmt, err)
		return fmt.Errorf("发布失败")
	} else if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		zap.S().Errorf("sql:%s, res.RowsAffected() return err: %v, RowsAffected:%d", sqlStmt, err, affected)
		return fmt.Errorf("发布失败")
	}

	zap.S().Infof("publish version %d for script %s", p.Version, p.Name)
	return nil
}
