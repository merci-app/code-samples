package request

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func StructToUrlValues(f interface{}) url.Values {
	data := url.Values{}

	val := reflect.ValueOf(f)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.Type().NumField(); i++ {

		tag := strings.Split(val.Type().Field(i).Tag.Get("json"), ",")[0]
		field := val.Field(i)

		if tag == "" {
			tag = val.Type().Field(i).Name
		}

		// campos não exportáveis serão ignorados por esta função.
		// Necessário pois o método `Interface()` dá panic em campos não exportáveis.
		if val.Type().Field(i).Name[0:1] == strings.ToLower(val.Type().Field(i).Name[0:1]) {
			continue
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		// ponteiros nil terão kind 'Invalid'
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data.Add(tag, strconv.FormatInt(field.Int(), 10))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data.Add(tag, strconv.FormatUint(field.Uint(), 10))

		case reflect.Float32, reflect.Float64:
			data.Add(tag, fmt.Sprint(field.Interface()))

		case reflect.String:
			data.Add(tag, field.String())

		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				data.Add(tag, fmt.Sprint(field.Index(i)))
			}

		case reflect.Struct:
			if field.Type().String() == "decimal.Decimal" {
				data.Add(tag, fmt.Sprint(field.Interface()))
			}

		case reflect.Bool:
			value := strconv.FormatBool(field.Bool())
			if value != "" {
				data.Add(tag, value)
			}
		}
	}

	return data
}
