package main

import (
	"errors"
	"fmt"
	"wrap-error/dao"
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
	fmt.Printf("main: query successfully result is %v\n", result)
}
