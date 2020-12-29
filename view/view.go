package view

import (
	"errors"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
	"io"
)

type ReadCloserReaderAt interface {
	io.ReadCloser
	io.ReaderAt
}

func NewView(r ReadCloserReaderAt, pageSize mysql_define.PAGE_SIZE) *View {
	return &View{
		buf:      ringbuffer.New(int(pageSize)),
		readPtr:  r,
		pageSize: pageSize,
	}
}

type View struct {
	buf *ringbuffer.RingBuffer
	// TODO need (*io.ReadCloser) ?
	readPtr  ReadCloserReaderAt
	pageSize mysql_define.PAGE_SIZE
}

var ErrFileEOF = errors.New("EOF")
var ErrFileRead = errors.New("read from file")
var ErrFileReadLength = errors.New("read length from file")

var ErrRingBufferWrite = errors.New("write to ringbuffer")
var ErrRingBufferRead = errors.New("read from ringbuffer")
var ErrRingBufferReadLength = errors.New("read lenght from ringbuffer")

var ErrPageParseFactory = errors.New("PageParseFactory error")

func (v *View) readAtPageNo(pageNo int) error {

	if !v.buf.IsEmpty() {
		v.buf.Reset()
	}

	var buffer = make([]byte, v.pageSize)

	length, err := v.readPtr.ReadAt(buffer, int64(pageNo*(int(v.pageSize))))
	if err != nil {
		log.Errorf("view.readNext; error:[%v]", err)
		if err == io.EOF {
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
		if err == io.EOF {
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
	// 1. pageNo is null, print all file type about page.
	// 2. if len(pageNo) is not zero, ready print file type contained by pageNo array.
	// 3. v.readPage(xxx) --> creating a object to parse file type (factory model)---->

	var (
		err         error
		fileAllPage IPageParse
	)

	if fileAllPage = new(PageParseFactory).Create(mysql_define.FIL_COMMON_HEADER_TAILER); fileAllPage == nil {
		return ErrPageParseFactory
	}

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
		fileAllPage.PrintPageType()
	}

}

func (v *View) ParsePage(pageNo int) error{
	// 1. check sum value, judge this page whether is correct.
	// 2. file header, file trailer
	// 3. file body {manger page, inode page, index page, freshly allocated page}

	var(
		err error
		fileAllPage IPageParse
		pageObject IPageParse
	)

	err = v.readAtPageNo(pageNo)
	if err != nil && err != ErrFileEOF {
		log.Errorf("view.readNext; error:[%v]", err)
		return err
	}
	if err == ErrFileEOF {
		//end
		return nil
	}

	if fileAllPage = new(PageParseFactory).Create(mysql_define.FIL_COMMON_HEADER_TAILER); fileAllPage == nil {
		return ErrPageParseFactory
	}
	fileAllPage.PageParseFILHeader(v.buf)


	log.Debugf("ready parse [%s]", mysql_define.StatusText(fileAllPage.GetFileType()))

	if pageObject = new(PageParseFactory).Create(fileAllPage.GetFileType()); pageObject == nil {
		return ErrPageParseFactory
	}

	pageObject.PageParseBody(v.buf, v.pageSize)
	pageObject.PrintPageType()

	return nil
}



