package helpers

import "reflect"

func PropertyExists(name string, data interface{}) bool {
    v := reflect.ValueOf(data)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    if v.Kind() != reflect.Struct {
        return false
    }
    return v.FieldByName(name).IsValid()
}