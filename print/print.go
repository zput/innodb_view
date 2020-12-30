package print

import (
	"github.com/modood/table"
	"reflect"
)

type foo func (t *PrintFormatT)

func SetPosition(p int)foo{
	return func(t *PrintFormatT){
		t.Position = p
	}
}
func SetName(n string)foo{
	return func(t *PrintFormatT){
		t.Name = n
	}
}
func SetValue(v interface{})foo{
	return func(t *PrintFormatT){
		t.Value = v
	}
}

func ConstructPrintFormatT(functions ...foo)(ret *PrintFormatT){
	ret = new(PrintFormatT)
	for index := range functions{
		functions[index](ret)
	}
	return
}

type PrintFormatT struct {
	Position  int `table:"position"`
	Name string `table:"name"`
	Value interface{} `table:"value"`
}

// Output to stdout
//table.Output(objects)
func PrintFun(objects []PrintFormatT)string{

	// Or just return table string and then do something
	return table.Table(objects)
}

func Translate(obj interface{}) []PrintFormatT {
	// Wrap the original in a reflect.Value
	var original = reflect.ValueOf(obj)
	var ret = make([]PrintFormatT, 0)

	translateRecursive(0, "", &ret, original)

	// Remove the reflection wrapper
	return ret
}

func translateRecursive(currentPosition int, prefix string, ret *[]PrintFormatT, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}

		translateRecursive(currentPosition, prefix, ret, originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()

		translateRecursive(currentPosition, prefix, ret, originalValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(currentPosition, prefix+"."+original.Type().Field(i).Tag.Get("self"), ret, original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(currentPosition, prefix, ret, original.Index(i))
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	//case reflect.String:
	//	*ret = append(*ret, *ConstructPrintFormatT(SetPosition(currentPosition), SetName(prefix), SetValue(original)))
	//	currentPosition += int(original.Type().Size())
	//
	//case reflect.Int:
	//	fallthrough
	//case reflect.Int8:
	//	fallthrough
	//case reflect.Int16:
	//	fallthrough
	//case reflect.Int32:
	//	fallthrough
	//case reflect.Int64:
	//	*ret = append(*ret, *ConstructPrintFormatT(SetPosition(currentPosition), SetName(prefix), SetValue(original)))
	//	currentPosition += int(original.Type().Size())

	// And everything else will simply be taken from the original
	default:
		*ret = append(*ret, *ConstructPrintFormatT(SetPosition(currentPosition), SetName(prefix), SetValue(original)))
		currentPosition += int(original.Type().Size())
	}

}
