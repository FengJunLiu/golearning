package main

import (
        "dao"
        "errors"
        "fmt"
)

func main() {
        result, err := dao.NewMySQL().Query("select name from users where id = ?", 1)
        if errors.Is(err, dao.ErrNoResult) {
                fmt.Printf("main: no result in table users where id=1\n")
                return
        }
        if err != nil {
                fmt.Printf("main: query failed %+v\n", err)
                return
        }
        fmt.Printf("main: query successfully result is %+v\n", result)
}

----------------------------------------------------------------------------------------
package dao

import (
        "database/sql"

        "errors"

        xerrors "github.com/pkg/errors"

        _ "github.com/go-sql-driver/mysql"
)

const (
        _DriverName     = "mysql"
        _DataSourceName = "root:root@tcp(127.0.0.1:3306)/test"
)

type MySQL struct {
        db  *sql.DB //mysql的连接句柄
        err error   //暂存error
}

var ErrNoResult = errors.New("dao: no rows in result set")

func NewMySQL() *MySQL {
        db, err := sql.Open(_DriverName, _DataSourceName)
        mysql := &MySQL{db: db}
        if err != nil {
                mysql.err = xerrors.Wrapf(err, "dao: NewMySQL failed")
                return mysql
        }
        if err := db.Ping(); err != nil {
                mysql.err = xerrors.Wrapf(err, "dao: NewMySQL failed")
                return mysql
        }
        return mysql
}

func (this *MySQL) Query(sqlStr string, vals ...interface{}) ([]map[string]string, error) {
        var result []map[string]string
        if this.err != nil {
                return result, this.err
        }
        rows, err := this.db.Query(sqlStr, vals...)
        if err != nil {
                return result, xerrors.Wrapf(err, "dao: Query failed")
        }
        defer rows.Close()

        cols, err := rows.Columns()
        if err != nil {
                return result, xerrors.Wrapf(err, "dao: Query failed")
        }
        col_len := len(cols)
        rawResult := make([][]byte, col_len)
        dest := make([]interface{}, col_len)
        for i, _ := range rawResult {
                dest[i] = &rawResult[i]
        }

        rowResult := make(map[string]string)
        for rows.Next() {
                err = rows.Scan(dest...)
                //sql.ErrNoRows
                if err != nil {
                        if err != sql.ErrNoRows {
                                return result, xerrors.Wrapf(err, "dao: Query failed")
                        }
                        return result, ErrNoResult
                }

                for i, raw := range rawResult {
                        key := cols[i]
                        if raw == nil {
                                rowResult[key] = ""
                        } else {
                                rowResult[key] = string(raw)
                        }
                }
                result = append(result, rowResult)

        }
        return result, nil
}
