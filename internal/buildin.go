package internal

import (
	"reflect"
)

var (
	mapTag map[reflect.Type]map[int]string = nil
)

func init() {
	mapTag = make(map[reflect.Type]map[int]string)
}

// InitTag 缓存Tag
func InitTag(t reflect.Type, tagName string) map[int]string {
	ma, mok := mapTag[t]
	if mok {
		return ma
	}
	ma = nil
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		tag := field.Tag
		if len(tag) > 0 {
			tv, ok := tag.Lookup(tagName)
			if ok {
				if ma == nil {
					ma = make(map[int]string)
					mapTag[t] = ma
				}
				ma[i] = tv
			}
		}
	}
	return ma
}
