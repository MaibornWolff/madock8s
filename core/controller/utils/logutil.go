package utils

import (
	"reflect"
)

type ResourceLogger struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

func SimplifyLog(d interface{}) ResourceLogger {
	v := reflect.ValueOf(d).Elem()
	r := ResourceLogger{
		Name:      v.FieldByName("Name").Interface().(string),
		Namespace: v.FieldByName("Namespace").Interface().(string),
		Labels:    v.FieldByName("Labels").Interface().(map[string]string),
	}
	return r
}
