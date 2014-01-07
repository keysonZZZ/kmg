package typeTransform

import (
	"reflect"
)

// copy fields value between two pointer of struct,
// it will copy value with same field name in two struct
// it ignore not share field name
// it ignore share field name with not same type
// it ignore field can not set (it is addressable and was not obtained by the use of unexported struct fields)
func StructFieldCopy(in interface{}, out interface{}) {
	vin := reflect.Indirect(reflect.ValueOf(in))
	tin := vin.Type()
	vout := reflect.Indirect(reflect.ValueOf(out))
	tout := vout.Type()
	if !vout.CanSet() {
		panic("[StructFieldCopy] out can not set,you have to passing a pointer.")
	}
	for i := 0; i < tin.NumField(); i++ {
		tinf := tin.Field(i)
		toutf, ok := tout.FieldByName(tinf.Name)
		if !ok {
			continue
		}
		if !tinf.Type.AssignableTo(toutf.Type) {
			continue
		}
		voutfv := vout.FieldByIndex(toutf.Index)
		if !voutfv.CanSet() {
			continue
		}
		voutfv.Set(vin.Field(i))
	}
}
