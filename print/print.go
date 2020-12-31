package print

import (
	"fmt"
	"github.com/modood/table"
	"reflect"
	"strings"
)

type foo func(t *PrintFormatT)

func SetPosition(p interface{}) foo {
	return func(t *PrintFormatT) {
		t.Position = p
	}
}
func SetPositionString(p string) foo {
	return func(t *PrintFormatT) {
		t.Position = p
	}
}
func SetName(n string) foo {
	return func(t *PrintFormatT) {
		t.Name = n
	}
}
func SetValue(v interface{}) foo {
	return func(t *PrintFormatT) {
		t.Value = v
	}
}

func ConstructPrintFormatT(functions ...foo) (ret *PrintFormatT) {
	ret = new(PrintFormatT)
	for index := range functions {
		functions[index](ret)
	}
	return
}

type PrintFormatT struct {
	Position interface{} `table:"POSITION"`
	Name     string      `table:"NAME"`
	Value    interface{} `table:"VALUE"`
}

type PrintDivideSign int

const (
	PrintDivideSignBlock = iota
	PrintDivideSignBlank
	PrintDivideSignTrailer
)

func NewPrintFormatT(sign PrintDivideSign, name string) *PrintFormatT {

	var headSplitSign = ""
	var trailSplitSign = ""
	var nameSign = ""

	switch sign {
	case PrintDivideSignBlock:
		headSplitSign = "**************"
		trailSplitSign = "**************"
		nameSign = fmt.Sprintf("**************%s**************", name)
	case PrintDivideSignBlank:

	case PrintDivideSignTrailer:
		headSplitSign = "N+8"
	}

	return ConstructPrintFormatT(SetPositionString(headSplitSign), SetName(nameSign), SetValue(trailSplitSign))
}

// Output to stdout
//table.Output(objects)
func PrintFun(objects []PrintFormatT) string {

	// Or just return table string and then do something
	return table.AsciiTable(objects)
}

func Translate(currentPosition interface{}, obj interface{}) []PrintFormatT {
	// Wrap the original in a reflect.Value
	var original = reflect.ValueOf(obj)
	var ret = make([]PrintFormatT, 0)

	translateRecursive(currentPosition, "", &ret, original)

	// Remove the reflection wrapper
	return ret
}

func translateRecursive(currentPositionPtr interface{}, prefix string, ret *[]PrintFormatT, original reflect.Value) {
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

		translateRecursive(currentPositionPtr, prefix, ret, originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()

		translateRecursive(currentPositionPtr, prefix, ret, originalValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(currentPositionPtr,
				strings.TrimPrefix(prefix+"."+snakeString(original.Type().Field(i).Tag.Get("self")), "."),
				ret, original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(currentPositionPtr, prefix+fmt.Sprintf("[%d]", i), ret, original.Index(i))
			*ret = append(*ret, *NewPrintFormatT(PrintDivideSignBlank, ""))
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
		var saveUpperFloorData interface{}
		tmpPointerValue := reflect.ValueOf(currentPositionPtr)
		switch tmpPointerValue.Elem().Type().Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			//(*currentPositionPtr).(int) += int(original.Type().Size())
			saveUpperFloorData = tmpPointerValue.Elem().Interface().(int)
			tmpPointerValue.Elem().SetInt(int64(saveUpperFloorData.(int) + int(original.Type().Size())))

		case reflect.String:
			//(*currentPositionPtr).(string) += fmt.Sprintf("+%d", int(original.Type().Size()))
			saveUpperFloorData = tmpPointerValue.Elem().Interface().(string)
			tmpPointerValue.Elem().SetString(saveUpperFloorData.(string) + fmt.Sprintf("+%d", int(original.Type().Size())))

		default:
			panic(fmt.Sprintf("print.translateRecursive; error; %v", tmpPointerValue.Elem().Type().Kind()))
		}

		*ret = append(*ret, *ConstructPrintFormatT(SetPosition(saveUpperFloorData), SetName(prefix), SetValue(original)))
	}

}

/**
 * 驼峰转蛇形 snake string
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @param s 需要转换的字符串
 * @return string
 **/
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j && (d!='D'&&s[i-1]!='I')&&(d!='S'&&s[i-1]!='L')&&(i>1&&d!='N'&&s[i-1]!='S'&&s[i-2]!='L'){
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}
