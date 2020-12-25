package main

import (
	flag "github.com/spf13/pflag"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/innodb_view/view"
	"os"
)

var (
	operatorTypePtr *string
	pageNumbers     *[]int
	pageSize        mysql_define.PAGE_SIZE
	filePathPtr     *string
)

const (
	OperatorType_SCAN  = "scan"
	OperatorType_PARSE = "parse"
)

/*
./programName
	--opertor_type(-oT) scan/parse
	--page_no_numbers(-pNNs)
 	--page_size(-pS) 16/32
	--file_path (-f)	wait parsing file
*/

func checkParam() {

	// https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/
	operatorTypePtr = flag.StringP("opertor_type", "t", "scan", "operator type:(scan/parse)")
	pageNumbers = flag.IntSliceP("page_numbers", "n", []int{0}, "page numbers: all page is [-1]; others is [0,1,...]")
	pageSizePtr := flag.IntP("page_size", "s", 16, "page size:(16/32 etc)")
	filePathPtr = flag.StringP("file_path", "f", "scan", "wait parsing file")
	debugModePtr := flag.BoolP("debug_mode", "d", false, "debug mode (default:false)")

	flag.Parse()

	// TODO judge file whether exist.

	switch *pageSizePtr {
	case 16:
		pageSize = mysql_define.PAGE_SIZE_16_K
	case 32:
		pageSize = mysql_define.PAGE_SIZE_32_K
	}

	if *debugModePtr == true {
		log.SetLevel(log.LevelDebug)
	}

}

func main() {

	checkParam()

	//http://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.1.html
	//func Open(name string) (*File, error)
	f, err := os.Open(*filePathPtr)
	defer f.Close() // need close file
	if err != nil {
		panic(err)
	}

	viewObject := view.NewView(f, pageSize)

	switch *operatorTypePtr {
	case OperatorType_SCAN:
		_ = viewObject.ScanPage()
	case OperatorType_PARSE:
		viewObject.ParsePage((*pageNumbers)[0])
	default:
		panic("parameter error; operator_type")
	}
}
