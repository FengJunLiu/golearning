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
		//暂存error
		mysql.err = xerrors.Wrapf(err, "dao: NewMySQL failed")
		return mysql
	}
	if err := db.Ping(); err != nil {
		//暂存error
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
		// 查询空值需要被业务层感知，业务层对空值一般情况下需要做降级处理
		// 在dao层定义sql.ErrNoRows对应的dao.ErrNoResult，上层只需引用dao层即可做空值的判定
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
