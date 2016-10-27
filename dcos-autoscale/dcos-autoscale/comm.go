package main

import  (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"

)

func GetOne(rows *sql.Rows) []string {

	fmt.Println("start")
    if rows == nil {
        return nil
    }  

    cols, err := rows.Columns() 
    CheckErr(err) 
    rawResult := make([][]byte, len(cols))
    result := make([]string, len(cols))
    dest := make([]interface{}, len(cols))
    for i, _ := range rawResult {
        dest[i] = &rawResult[i]
    }
    
    if rows.Next() {
        err = rows.Scan(dest...)
        CheckErr(err)
        for i, raw := range rawResult {
            if raw == nil {
                result[i] = ""
            } else {
                result[i] = string(raw)
                fmt.Println(result[i])
            }
        }
        
        //fmt.Printf("%#v\n", result)
        
        //break
    } else {
        return nil
    }
    
    _=err
    return result
}