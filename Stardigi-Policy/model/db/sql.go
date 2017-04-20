package db

func NewSql(name string) Sqler {

	sqler := GetDBsqker(name)

	return sqler
}
