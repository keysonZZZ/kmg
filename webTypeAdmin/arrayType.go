package webTypeAdmin

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
)

//path -> array index(number)
type arrayType struct {
	commonType
	elemType typeInterface
}

func (t *arrayType) init() {
	if t.elemType != nil {
		return
	}
	t.elemType = t.ctx.mustNewTypeFromReflect(t.getReflectType().Elem())
}
func (t *arrayType) Html(v reflect.Value) template.HTML {
	t.init()
	type templateRow struct {
		Path  int
		Index int
		Html  template.HTML
	}
	var templateData []templateRow
	len := v.Len()
	for i := 0; i < len; i++ {
		templateData = append(templateData, templateRow{
			Path:  i,
			Index: i,
			Html:  t.elemType.Html(v.Index(i)),
		})
	}
	return theTemplate.MustExecuteNameToHtml("Array", templateData)
}

func (t *arrayType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	t.init()
	return arrayGetSubValueByString(v, k)
}

func arrayGetSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	i, err := arrayParseKey(v, k)
	if err != nil {
		return reflect.Value{}, nil
	}
	return v.Index(i), nil
}
func arrayParseKey(v reflect.Value, k string) (int, error) {
	i64, err := strconv.ParseInt(k, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[arrayParseKey] index is not int k:%s", k)
	}
	i := int(i64)
	if i >= v.Len() || i < 0 {
		return 0, fmt.Errorf("[arrayParseKey] index is not of range k:%s,len:%d", k, v.Len())
	}
	return i, nil
}
