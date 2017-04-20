package utils

import (
	"bytes"
	 "fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StringJoin(str ...string) string {

	var buffer bytes.Buffer
	for _, s := range str {
		buffer.WriteString(s)
	}

	return buffer.String()
}

func InterfaceJoin(in ...interface{}) string {

	var buffer bytes.Buffer
	for _, s := range in {
		switch s.(type) {
		case int:
			buffer.WriteString(strconv.Itoa(s.(int)))
		case string:
			buffer.WriteString(s.(string))
		}
	}
	return buffer.String()
}

func StructJoin(object interface{}) string {

	_, field := ParaseInterface(object)
	str := InterfaceJoin(field...)
	return str
}

func ParaseUrlParam(url string, params map[string]string) string {

	// /v2/apps/:app_id map[app_id:/super/admin]
	for k, v := range params {
		//判断v是不是以/开始
		if StartWith(v, `/`) {
			k = StringJoin(`/:`, k)
		}
		url = strings.Replace(url, k, v, -1)
	}

	return url
}

func TagParase(object interface{}, arg ...string) map[string]string {

	tags := make(map[string]string)
	typ := reflect.TypeOf(object)
	field := typ.Elem().Field(0)
	for _, a := range arg {
		tags[a] = field.Tag.Get(a)
	}

	return tags
}

/**
 * the function parase string like 'AppScaleRule' to 'app_scale_rule'
 */
func StringParaseWith_(str string) string {

	var buffer bytes.Buffer
	for i, s := range str {
		if s >= 65 && s <= 90 && i != 0 {
			buffer.WriteByte('_')
			buffer.WriteRune(s)
			continue
		}
		buffer.WriteRune(s)
	}
	fmt.Println("--------------name-----",strings.ToLower(buffer.String()))
	return strings.ToLower(buffer.String())
}

/**
 * the function declare to get field of a given struct
 */
func NumberOfContainer(object interface{}) int {

	typ := reflect.TypeOf(object)
	numField := typ.Elem().NumField()
	return numField
}

func NameParaseWithReflect(object interface{}) string {

	typ := reflect.TypeOf(object)
	// 判断typ的类型是指针类型，还是数据类型
	typRmPtr := VerifyType(typ)
	fmt.Println("-----typeRmPtr",typRmPtr.Name())
	return typRmPtr.Name()
}

func ParaseInterfaceWithField(container interface{}, arg []string) interface{} {

	// var tmp int64
	typ := reflect.ValueOf(container)
	elem := typ.Elem()
	numField := elem.NumField()
	for i := 0; i < numField; i++ {
		switch elem.Field(i).Kind() {
		case reflect.String:
			elem.Field(i).SetString(arg[i])
		case reflect.Int:
			tmp, _ := strconv.ParseInt(arg[i], 10, 64)
			elem.Field(i).SetInt(tmp)
		case reflect.Float64:
			tmp, _ := strconv.ParseFloat(arg[i], 64)
			elem.Field(i).SetFloat(tmp)
		}
	}

	return container
}

func UrlBase64(urls string) string {

	return url.QueryEscape(urls)
}

func StartWith(s string, sep string) bool {

	if i := strings.Index(s, sep); i == 0 {
		return true
	}
	return false
}

func ParaseInterface(object interface{}) (string, []interface{}) {
	var value []interface{}
	var buffer bytes.Buffer
	valueof := reflect.ValueOf(object)
	numField := valueof.Elem().NumField()

	for i := 0; i < numField; i++ {
		value = append(value, valueof.Elem().Field(i).Interface())
		name := StringParaseWith_(valueof.Elem().Type().Field(i).Name)
		str := StringJoin(name, "=?")
		buffer.WriteString(str)
		if i < numField-1 {
			buffer.WriteString(",")
		}
	}

	return buffer.String(), value
}

func CurrentimeString() string {

	timestamp := time.Now().Unix()
	timeString := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

	return timeString
}
