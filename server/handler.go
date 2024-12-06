package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
)

type MysqlHandler struct {
}

func (m *MysqlHandler) UseDB(dbName string) error {
	return nil
}

func (m *MysqlHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	return nil, nil
}

func (m *MysqlHandler) HandleStmtClose(context interface{}) error {
	return nil
}

func (m *MysqlHandler) HandleQuery(query string) (*mysql.Result, error) {
	res, err := mysql.BuildSimpleBinaryResultset([]string{"version"}, [][]any{
		{"8.0.12"},
	})
	if err != nil {
		return nil, err
	}
	return &mysql.Result{
		Status:       0,
		Warnings:     0,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    res,
	}, nil
}

func (m *MysqlHandler) HandleStmtPrepare(query string) (params int, columns int, context interface{}, err error) {
	ss := strings.Split(query, " ")
	switch strings.ToLower(ss[0]) {
	case "select", "insert", "update", "delete", "replace":
		params = strings.Count(query, "?")
	default:
		err = fmt.Errorf("invalid prepare %s", query)
	}
	return params, columns, nil, err
}

func (m *MysqlHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	res, err := mysql.BuildSimpleBinaryResultset([]string{"id", "created_at", "updated_at"}, [][]any{
		{1, time.Now(), time.Now()},
	})
	if err != nil {
		return nil, err
	}
	return &mysql.Result{
		Status:       0,
		Warnings:     0,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    res,
	}, nil

}

func (m *MysqlHandler) HandleOtherCommand(cmd byte, data []byte) error {
	return nil
}
