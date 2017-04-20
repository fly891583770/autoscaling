package db

type Sqler interface {
	QueryTable(object interface{}) QuerySet
	Insert(object interface{})
}

type QuerySet interface {
	Filter(col string, arg interface{}) QuerySet
	One(container interface{}, cols ...string) error
	All(container interface{}, cols ...string) error
	ReadRows(container interface{}) ([][]string, error)
}
