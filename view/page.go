package view

import (
	"encoding/json"
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
)

type ListBaseNode struct{
	FlstLen uint32
	First PageNoANDOffset
	Last PageNoANDOffset
}

type PageNoANDOffset struct{
	PageNo uint32
	Offset uint16
}

//---------------------------------------------------------------------//

type IPageParse interface {
	PageParseFILHeader(buffer *ringbuffer.RingBuffer) error
	PageParseFILTailer(buffer *ringbuffer.RingBuffer) error
	PageParseBody(buffer *ringbuffer.RingBuffer) error
	PrintPageType()error

	GetFileType()mysql_define.T_FIL_PAGE_TYPE
}

type PageParseFactory struct{}

func (f *PageParseFactory) Create(pageType mysql_define.T_FIL_PAGE_TYPE) IPageParse {
	switch pageType {
	case mysql_define.FIL_COMMON_HEADER_TAILER:
		return new(FileAllPage)
	case mysql_define.FIL_PAGE_TYPE_FSP_HDR:
		return new(FspHeaderPage)

	case mysql_define.FIL_PAGE_INODE:

	case mysql_define.FIL_PAGE_INDEX:

	case mysql_define.FIL_PAGE_TYPE_XDES:

	case mysql_define.FIL_PAGE_TYPE_ALLOCATED:

	}
	return nil
}

type FileAllPage struct {
	checksum                    uint32
	Offset                      uint32
	previousPage                uint32
	nextPage                    uint32
	lsnForLastPageModeification uint64
	pageType                    uint16
	flushLSN                    uint64 //(0 except space0 page0)
	spaceID                     uint32

	//FIL_PAGE_SPACE_OR_CHKSUM = 0
	//FIL_PAGE_OFFSET = 4
	//FIL_PAGE_PREV = 8
	//FIL_PAGE_NEXT = 12
	//FIL_PAGE_LSN = 16
	//FIL_PAGE_TYPE = 24
	//FIL_PAGE_FILE_FLUSH_LSN = 26
	//FIL_NULL = 0xFFFFFFFF /*no PAGE_NEXT or PAGE_PREV */
	//FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID = 34

	oldStyleChecksum uint32
	lsn              uint32
	// FIL_PAGE_END_LSN_OLD_CHKSUM 8 byte

}

func (fap *FileAllPage) GetFileType()mysql_define.T_FIL_PAGE_TYPE{
	return mysql_define.T_FIL_PAGE_TYPE(fap.pageType)
}

func (fap *FileAllPage) printPageType() error {
	log.Debugf("page type value:%d", fap.pageType)

	pageNumber := fap.Offset

	switch mysql_define.T_FIL_PAGE_TYPE(fap.pageType) {

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

func (fap *FileAllPage) PrintPageType() error {

	if err := fap.printPageType(); err != nil{
		return err
	}
	return nil
}

func (fap *FileAllPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {

	var isUsingExplore = true

	buffer.ExploreBegin()

	// TODO optimize this handler error
	// FIL_PAGE_OFFSET
	fap.checksum = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_OFFSET); err != nil {
		log.Error(err)
		return err
	}

	fap.Offset = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_PREV - mysql_define.FIL_PAGE_OFFSET); err != nil {
		log.Error(err)
		return err
	}

	fap.previousPage = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_NEXT - mysql_define.FIL_PAGE_PREV); err != nil {
		log.Error(err)
		return err
	}

	fap.nextPage = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_LSN - mysql_define.FIL_PAGE_NEXT); err != nil {
		log.Error(err)
		return err
	}

	fap.lsnForLastPageModeification = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_TYPE - mysql_define.FIL_PAGE_LSN); err != nil {
		log.Error(err)
		return err
	}

	fap.pageType = buffer.PeekUint16(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_FILE_FLUSH_LSN - mysql_define.FIL_PAGE_TYPE); err != nil {
		log.Error(err)
		return err
	}

	fap.flushLSN = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID - mysql_define.FIL_PAGE_FILE_FLUSH_LSN); err != nil {
		log.Error(err)
		return err
	}

	fap.spaceID = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_DATA - mysql_define.FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID); err != nil {
		log.Error(err)
		return err
	}

	buffer.ExploreBreak()


	return nil
}

func (fap *FileAllPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer) error {

	var isUsingExplore = true

	buffer.ExploreBegin()

	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_TRAILER_BEFORE_SIZE); err != nil {
		log.Error(err)
		return err
	}
	fap.oldStyleChecksum = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_END_LSN_OLD_CHKSUM/2); err != nil {
		log.Error(err)
		return err
	}

	fap.lsn = buffer.PeekUint32(isUsingExplore)

	buffer.ExploreBreak()

	return nil
}

func (fap *FileAllPage) PageParseBody(buffer *ringbuffer.RingBuffer) error {
	return nil
}

// ----------------- FspHeaderPage ------------------------------------//
type FspHeaderPage struct{
	FileAllPage
	FspSpaceID uint32
	FspNotUsed uint32
	FspSize uint32
	FspFreeLimit uint32
	FspSpaceFlags uint32
	FspFragNUsed uint32

	FspFree *ListBaseNode
	FspFreeFrag *ListBaseNode
	FspFullFrag *ListBaseNode

	FspSegID uint64

	FspSegInodesFull *ListBaseNode
	FspSegInodesFree *ListBaseNode
}

func (fhp *FspHeaderPage) GetFileType()mysql_define.T_FIL_PAGE_TYPE{
	return mysql_define.T_FIL_PAGE_TYPE(fhp.FileAllPage.pageType)
}

func (fhp *FspHeaderPage) printPageType() error {
	prettyFormat, err := json.MarshalIndent(fhp, "", "    ")
	if err != nil{
		return err
	}
	fmt.Printf("%s", string(prettyFormat))
	return nil
}

func (fhp *FspHeaderPage) PrintPageType() error {
	fhp.printPageType()

	if err := fhp.FileAllPage.PrintPageType(); err != nil{
		log.Error(err)
		return err
	}


	return nil
}
func (fhp *FspHeaderPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {
	if err := fhp.FileAllPage.PageParseFILHeader(buffer); err != nil{
		return err
	}

	return nil
}
func (fhp *FspHeaderPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer) error {
	if err := fhp.FileAllPage.PageParseFILTailer(buffer); err != nil{
		return err
	}

	return nil
}
func (fhp *FspHeaderPage) PageParseBody(buffer *ringbuffer.RingBuffer) error {

	var isUsingExplore = true

	buffer.ExploreBegin()

	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_DATA); err != nil {
		log.Error(err)
		return err
	}
	fhp.FspSpaceID = buffer.PeekUint32(isUsingExplore)

	log.Debugf("FSP spaceID[%d]", fhp.FspSpaceID)

	if err := buffer.ExploreRetrieve(mysql_define.FSP_NOT_USED); err != nil {
		log.Error(err)
		return err
	}
	fhp.FspNotUsed = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_SIZE - mysql_define.FSP_NOT_USED); err != nil {
		log.Error(err)
		return err
	}
	fhp.FspSize = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_FREE_LIMIT - mysql_define.FSP_SIZE); err != nil {
		log.Error(err)
		return err
	}

	fhp.FspFreeLimit = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_SPACE_FLAGS - mysql_define.FSP_FREE_LIMIT); err != nil {
		log.Error(err)
		return err
	}
	fhp.FspSpaceFlags = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_FRAG_N_USED - mysql_define.FSP_SPACE_FLAGS); err != nil {
		log.Error(err)
		return err
	}
	fhp.FspFragNUsed = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_FREE - mysql_define.FSP_FRAG_N_USED); err != nil {
		log.Error(err)
		return err
	}

	var err error
	if fhp.FspFree, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return err
	}
	if fhp.FspFreeFrag, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return err
	}
	if fhp.FspFullFrag, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return err
	}

	fhp.FspSegID = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_SEG_INODES_FULL - mysql_define.FSP_SEG_ID); err != nil {
		log.Error(err)
		return err
	}

	if fhp.FspSegInodesFull, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return err
	}
	if fhp.FspSegInodesFree, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return err
	}

	buffer.ExploreBreak()

	return nil
}





type InodePage struct{}

type IndexPage struct{}



func getListBaseNode(buffer *ringbuffer.RingBuffer)(*ListBaseNode, error){

	var isUsingExplore = true
	var listBaseNode ListBaseNode

	listBaseNode.FlstLen = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.PRV_PAGE_NODE); err != nil {
		log.Error(err)
		return nil, err
	}
	listBaseNode.First.PageNo = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.PRV_OFFSET - mysql_define.PRV_PAGE_NODE); err != nil {
		log.Error(err)
		return nil, err
	}

	listBaseNode.First.Offset = buffer.PeekUint16(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.NEXT_PAGE_NODE - mysql_define.PRV_OFFSET); err != nil {
		log.Error(err)
		return nil, err
	}

	listBaseNode.Last.PageNo = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.NEXT_OFFSET - mysql_define.NEXT_PAGE_NODE); err != nil {
		log.Error(err)
		return nil, err
	}

	listBaseNode.Last.Offset = buffer.PeekUint16(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.LIST_BASE_NODE_SIZE - mysql_define.NEXT_OFFSET); err != nil {
		log.Error(err)
		return nil, err
	}

	return &listBaseNode, nil
}
