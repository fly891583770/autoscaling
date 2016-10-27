package  main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
//	"strings"
)

func ModelinterFace(args ...interface{}) int64 {

	db,err:= sql.Open("mysql",macconfig.Datasource)
	CheckErr(err)
	stmt, err := db.Prepare(args[0].(string))
	CheckErr(err)

	dest := make([]interface{},0)

	for k,v :=range args{
		if k > 0 {
			dest = append(dest,v)
		}
	}

	res, err := stmt.Exec(dest...)
	CheckErr(err)
	affect, err := res.RowsAffected()
	CheckErr(err)
	fmt.Println("affect:", affect)
	db.Close()
	return affect
	
}

func SerachInfo(args ...interface{}) []string{

	db,err:= sql.Open("mysql",macconfig.Datasource)
	CheckErr(err)	
	stmt, err := db.Prepare(args[0].(string))
	CheckErr(err)

	dest := make([]interface{},0)
	for k,v :=range args{
		if k > 0 {
			dest = append(dest,v)
		}
	}
	rows, err := stmt.Query(dest...)
	CheckErr(err)
	result := GetOne(rows)
	stmt.Close()
	return result
}


func SerachAppid(args ...interface{}) []string{

	var results []string
	var result string
	db,err:= sql.Open("mysql",macconfig.Datasource)
	CheckErr(err)

	dest := make([]interface{},0)
	for k,v :=range args{
		if k > 0 {
			dest = append(dest,v)
		}
	}

	rows, err := db.Query(args[0].(string),dest...)
	CheckErr(err)
	for index := 0; rows.Next(); index++ {
		err = rows.Scan(&result)
		CheckErr(err)
		results = append(results, result)
	}
	db.Close()
	return results
}