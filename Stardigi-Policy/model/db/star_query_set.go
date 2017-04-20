package db

import (
	"Stardigi-Policy/utils"
	"fmt"
	"reflect"
)

type StarQuerySet struct {
	DbClient *DBclient
	Cond     bool
}

func NewStarQuerySet(dbClient *DBclient) QuerySet {

	return &StarQuerySet{
		DbClient: dbClient,
		Cond:     false,
	}
}

/**
 * "SELECT * from ", table, " where "
 */
func (star StarQuerySet) Filter(col string, arg interface{}) QuerySet {

	var strCon string
	strSql := star.DbClient.StrSql
	//判断之前是否有条件存在
	if !star.Cond {
		strSql = utils.StringJoin(strSql, " where ")
		argWithQuota := utils.StringJoin(`"`, arg.(string), `"`)
		strCon = utils.StringJoin(strSql, col, "=", argWithQuota)
		star.Cond = true
	} else {
		strSql = utils.StringJoin(strSql, " and ")
		argWithQuota := utils.StringJoin(`"`, arg.(string), `"`)
		strCon = utils.StringJoin(strSql, col, "=", argWithQuota)
		star.Cond = true
	}

	// "SELECT * from ", table, " where col = arg AND "

	// 保存查询语句
	star.DbClient.StrSql = strCon
	return star
}

func (star StarQuerySet) All(container interface{}, cols ...string) error {

	// 获取mysql对应的bean对象
	object := star.DbClient.object
	// 获取所有的结果
	rows, err := star.ReadRows(object)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(container)
	ind := reflect.Indirect(val)
	// slice := ind
	// fmt.Println("======ind.type=====", ind.Type().Elem())
	slice := reflect.New(ind.Type()).Elem()
	// fmt.Println("======ind.type=====", slice.Type().Elem())
	// 将数据填充到容器
	for _, r := range rows {
		con := utils.ParaseInterfaceWithField(object, r)
		slice = reflect.Append(slice, reflect.ValueOf(con).Elem())

	}
	ind.Set(slice)

	return nil
}

/**
 * the function return single row
 * if mutiple rows return from db,the function will return the first one
 */
func (star StarQuerySet) One(container interface{}, cols ...string) error {

	rows, err := star.ReadRows(container)
	if err != nil {
		return err
	}

	for _, r := range rows {
		// 填充数据到容器
		utils.ParaseInterfaceWithField(container, r)
		break
	}
	return nil
}

func (star StarQuerySet) ReadRows(object interface{}) ([][]string, error) {
	// 获取dbclient
	db := star.DbClient.Db
	// 获取sql语句
	strChan := star.DbClient.StrSql
	numContainer := utils.NumberOfContainer(object)
	rowContainer := make([][]string, 0)
	rows, err := db.Query(strChan)
	if err != nil {
		return rowContainer, err
	}
	defer rows.Close()
	col := make([]interface{}, numContainer)
	for i := range col {
		var co interface{}
		col[i] = &co
	}

	for rows.Next() {
		err = rows.Scan(col...)
		if err != nil {
			fmt.Println("======scan发送错误======", err)
			return rowContainer, err
		}
		colContainer := make([]string, 0)
		for _, c := range col {
			val := reflect.Indirect(reflect.ValueOf(c)).Interface()
			switch v := val.(type) {
			case []byte:
				colContainer = append(colContainer, string(v))
			}
		}
		rowContainer = append(rowContainer, colContainer)
	}

	return rowContainer, nil
}
