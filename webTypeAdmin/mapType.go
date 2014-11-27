package webTypeAdmin

import (
	"html/template"
	"reflect"

	"github.com/bronze1man/kmg/kmgType"
)

//path -> key(Key type)
//key can be bool,string,stringEnum,int,float,
type mapType struct {
	kmgType.MapType
	elemAdminType adminType
	ctx           *context
}

//key is not string?
func (t *mapType) init() (err error) {
	if t.elemAdminType != nil {
		return
	}
	if err = t.MapType.Init(); err != nil {
		return
	}
	t.elemAdminType, err = t.ctx.typeOfFromKmgType(t.MapType.ElemType)
	if err != nil {
		return
	}
	return nil
}
func (t *mapType) HtmlView(v reflect.Value) (html template.HTML, err error) {
	if err = t.init(); err != nil {
		return
	}
	type templateRow struct {
		Path string
		Key  string
		Html template.HTML
	}
	var templateData []templateRow
	for _, key := range v.MapKeys() {
		keyS := t.KeyStringConverter.ToString(key)
		val := v.MapIndex(key)
		if html, err = t.elemAdminType.HtmlView(val); err != nil {
			return
		}
		templateData = append(templateData, templateRow{
			Path: keyS,
			Key:  keyS,
			Html: html,
		})
	}
	return theTemplate.ExecuteNameToHtml("Map", templateData)
}
