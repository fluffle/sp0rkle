package db

import "reflect"

// slicePtr makes working with a reflected pointer-to-slice less arduous.
type slicePtr struct {
	// pv is the reflect.Value of the pointer to the slice
	pv reflect.Value
	// sv is the reflect.Value of the slice itself
	sv reflect.Value
	// et is the element type the slice contains
	et reflect.Type
}

func newSlicePtr(value interface{}) *slicePtr {
	pv := reflect.ValueOf(value)
	if pv.Kind() != reflect.Ptr || pv.Elem().Kind() != reflect.Slice {
		panic("provided value is not a pointer-to-slice")
	}
	return &slicePtr{
		pv: pv,
		sv: pv.Elem(),
		et: pv.Elem().Type().Elem(),
	}
}

func (sp *slicePtr) newElem() reflect.Value {
	return reflect.New(sp.et).Elem()
}

func (sp *slicePtr) newStruct() reflect.Value {
	et := sp.et
	for et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	return reflect.New(et).Elem()
}

func (sp *slicePtr) appendElem(ev reflect.Value) {
	sp.sv = reflect.Append(sp.sv, ev)
	// Append may have returned a new slice so ensure pointer points to it.
	sp.pv.Elem().Set(sp.sv)
}

// ... I want a pony and this might just give me one.
func (sp *slicePtr) ponyElem() interface{} {
	ev := sp.newElem()
	sp.appendElem(ev)
	return sp.sv.Index(sp.len() - 1).Addr().Interface()
}

func (sp *slicePtr) len() int {
	return sp.sv.Len()
}
