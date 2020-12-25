package view

import (
	"errors"
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
	"io"
)

func NewView(r io.ReadCloser, pageSize mysql_define.PAGE_SIZE) *View {
	return &View{
		buf:      ringbuffer.New(int(pageSize)),
		readPtr:  r,
		pageSize: pageSize,
	}
}

type View struct {
	buf *ringbuffer.RingBuffer
	// TODO need (*io.ReadCloser) ?
	readPtr  io.ReadCloser
	pageSize mysql_define.PAGE_SIZE
}

var ErrFileEOF = errors.New("EOF")
var ErrFileRead = errors.New("read from file")
var ErrFileReadLength = errors.New("read length from file")

var ErrRingBufferWrite = errors.New("write to ringbuffer")
var ErrRingBufferRead = errors.New("read from ringbuffer")
var ErrRingBufferReadLength = errors.New("read lenght from ringbuffer")

//	func (v *View) readPage(pageNo int) {
func (v *View) readNext() error {
	// 1. clear ringBuffer
	// 2. read Page to ringBuffer, size is v.pageSize

	if !v.buf.IsEmpty() {
		v.buf.Reset()
	}

	var buffer = make([]byte, v.pageSize)

	length, err := v.readPtr.Read(buffer)
	if err != nil {
		log.Errorf("view.readNext; error:[%v]", err)
		if err == io.EOF{
			return ErrFileEOF
		}
		return ErrFileRead
	}
	if length != int(v.pageSize) {
		panic("can't get complete page size")
	}

	length, err = v.buf.Write(buffer)
	// TODO optimize handled error
	if err != nil {
		panic(err)
	}
	if length != int(v.pageSize) {
		panic("inner write error")
	}

	return nil
}

func (v *View) ScanPage() error {
	// TODO
	// 1. pageNo is null, print all file type about page.
	// 2. if len(pageNo) is not zero, ready print file type contained by pageNo array.
	// 3. v.readPage(xxx) --> creating a object to parse file type (factory model)---->

	var (
		err         error
		fileAllPage = new(PageParseFactory).Create(mysql_define.FIL_COMMON_HEADER_TAILER)
	)

	for {
		err = v.readNext()
		if err != nil && err != ErrFileEOF {
			log.Errorf("view.readNext; error:[%v]", err)
			return err
		}
		if err == ErrFileEOF {
			//end
			return nil
		}

		fileAllPage.PageParseFILHeader(v.buf)
	}

}

func (v *View) ParsePage(pageNo ...int) {
	// TODO
	// 1. check sum value, judge this page whether is correct.
	// 2. file header, file trailer
	// 3. file body {manger page, inode page, index page, freshly allocated page}

}

type IPageParse interface {
	PageParseFILHeader(buffer *ringbuffer.RingBuffer) error
	PageParseFILTailer(buffer *ringbuffer.RingBuffer) error
	PageParseBody(buffer *ringbuffer.RingBuffer) error
}

type PageParseFactory struct{}

func (f *PageParseFactory) Create(pageType mysql_define.T_FIL_PAGE_TYPE) IPageParse {
	switch pageType {
	case mysql_define.FIL_COMMON_HEADER_TAILER:
		return new(FileAllPage)
	case mysql_define.FIL_PAGE_TYPE_FSP_HDR:

	case mysql_define.FIL_PAGE_INODE:

	case mysql_define.FIL_PAGE_INDEX:

	case mysql_define.FIL_PAGE_TYPE_XDES:

	case mysql_define.FIL_PAGE_TYPE_ALLOCATED:

	}
	return nil
}

type FileAllPage struct{}

func (fap *FileAllPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {

	// FIL_PAGE_OFFSET
	buffer.Retrieve(mysql_define.FIL_PAGE_OFFSET)
	pageNumber := buffer.PeekUint32()

	buffer.Retrieve(mysql_define.FIL_PAGE_TYPE - mysql_define.FIL_PAGE_OFFSET)
	pageTypeValue := buffer.PeekUint16()

	log.Debugf("page type value:%d", pageTypeValue)

	switch mysql_define.T_FIL_PAGE_TYPE(pageTypeValue) {
	case mysql_define.FIL_PAGE_TYPE_FSP_HDR:
		fmt.Printf("page number:[%d]; page type:[%s]\n", pageNumber, mysql_define.StatusText(mysql_define.FIL_PAGE_TYPE_FSP_HDR))

	case mysql_define.FIL_PAGE_INODE:
		fmt.Printf("page number:[%d]; page type:[%s]\n", pageNumber, mysql_define.StatusText(mysql_define.FIL_PAGE_INODE))

	case mysql_define.FIL_PAGE_INDEX:
		fmt.Printf("page number:[%d]; page type:[%s]\n", pageNumber, mysql_define.StatusText(mysql_define.FIL_PAGE_INDEX))

	case mysql_define.FIL_PAGE_TYPE_ALLOCATED:
		fmt.Printf("page number:[%d]; page type:[%s]\n", pageNumber, mysql_define.StatusText(mysql_define.FIL_PAGE_TYPE_ALLOCATED))

	case mysql_define.FIL_PAGE_TYPE_XDES:
		fmt.Printf("page number:[%d]; page type:[%s]\n", pageNumber, mysql_define.StatusText(mysql_define.FIL_PAGE_TYPE_XDES))

	case mysql_define.FIL_PAGE_IBUF_BITMAP:
		fmt.Printf("page number:[%d]; page type:[%s]\n", pageNumber, mysql_define.StatusText(mysql_define.FIL_PAGE_IBUF_BITMAP))
	}
	return nil
}

func (fap *FileAllPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer) error {
	return nil
}

func (fap *FileAllPage) PageParseBody(buffer *ringbuffer.RingBuffer) error {
	return nil
}

type FspHeaderPage struct{}

type InodePage struct{}

type IndexPage struct{}
