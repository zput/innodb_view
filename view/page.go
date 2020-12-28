package view

import (
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
)

type ListBaseNode struct{
	FlstLen uint32
	First *PageNoANDOffset
	Last *PageNoANDOffset
}

type ListNode struct{
	First *PageNoANDOffset
	Last *PageNoANDOffset
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
		return new(INodePage)

	case mysql_define.FIL_PAGE_INDEX:

	case mysql_define.FIL_PAGE_TYPE_XDES:

	case mysql_define.FIL_PAGE_TYPE_ALLOCATED:

	}
	return nil
}

// ----------------- Fil header trailer ------------------------------------//

type FileAllPage struct {
	checksum                    uint32 `yaml:"checksum"`
	Offset                      uint32 `yaml:"Offset"`
	previousPage                uint32 `yaml:"previousPage"`
	nextPage                    uint32 `yaml:"nextPage"`
	lsnForLastPageModeification uint64 `yaml:"lsnForLastPageModeification"`
	pageType                    uint16 `yaml:"pageType"`
	flushLSN                    uint64 `yaml:"flushLSN"` //(0 except space0 page0) `yaml:"name"`
	spaceID                     uint32 `yaml:"spaceID"`

	//FIL_PAGE_SPACE_OR_CHKSUM = 0
	//FIL_PAGE_OFFSET = 4
	//FIL_PAGE_PREV = 8
	//FIL_PAGE_NEXT = 12
	//FIL_PAGE_LSN = 16
	//FIL_PAGE_TYPE = 24
	//FIL_PAGE_FILE_FLUSH_LSN = 26
	//FIL_NULL = 0xFFFFFFFF /*no PAGE_NEXT or PAGE_PREV */
	//FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID = 34

	oldStyleChecksum uint32 `yaml:"oldStyleChecksum"`
	lsn              uint32 `yaml:"lsn"`
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

// ----------------- common function ------------------------------------//

func getListBaseNode(buffer *ringbuffer.RingBuffer)(*ListBaseNode, error){

	var isUsingExplore = true
	var listBaseNode ListBaseNode
	var err error

	listBaseNode.FlstLen = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.PRV_PAGE_NODE); err != nil {
		log.Error(err)
		return nil, err
	}

	listBaseNode.First, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	listBaseNode.Last, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	return &listBaseNode, nil
}

func getPageNoANDOffset(buffer *ringbuffer.RingBuffer, isUsingExplore bool)(*PageNoANDOffset, error){

	var listNode PageNoANDOffset

	listNode.PageNo = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.PRV_OFFSET - mysql_define.PRV_PAGE_NODE); err != nil {
		log.Error(err)
		return nil, err
	}

	listNode.Offset = buffer.PeekUint16(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.NEXT_PAGE_NODE - mysql_define.PRV_OFFSET); err != nil {
		log.Error(err)
		return nil, err
	}

	return &listNode, nil
}

func getINodeEntry(buffer *ringbuffer.RingBuffer, isUsingExplore bool)(INodeEntry, error){

	var iNodeEntry INodeEntry
	var err error

	iNodeEntry.FSegID = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSEG_ID); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}

	iNodeEntry.FSegNotFullNUsed = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSEG_NOT_FULL_N_USED); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}

	if iNodeEntry.FSegFree, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return iNodeEntry, err
	}
	if iNodeEntry.FSegNotFull, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return iNodeEntry, err
	}
	if iNodeEntry.FSegFull, err = getListBaseNode(buffer); err != nil{
		log.Error(err)
		return iNodeEntry, err
	}

	iNodeEntry.FSegMagicN = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSEG_MAGIC_N); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}

	for i:=0; i<32;i++{
		var fSegFragTemp uint32

		if fSegFragTemp, err = getFSegFragArr(buffer, isUsingExplore); err != nil{
			log.Errorf("index[%d]; error[%v]", i, err)
			return iNodeEntry, err
		}
		if fSegFragTemp == 0xffffffff{
			log.Debug("all FSEG_FRAG_ARR_I object have showed in this segment object on INode page")
			if err = buffer.ExploreRetrieve(mysql_define.FSEG_FRAG_ARR_I*(31-i)); err != nil {
				log.Error(err)
				return iNodeEntry, err
			}
			break
		}

		iNodeEntry.FSegFragSlice = append(iNodeEntry.FSegFragSlice, fSegFragTemp)
	}

	return iNodeEntry, nil
}

func getFSegFragArr(buffer *ringbuffer.RingBuffer, isUsingExplore bool)(fSegFrag uint32, err error){
	fSegFrag = buffer.PeekUint32(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.FSEG_FRAG_ARR_I); err != nil {
		log.Error(err)
		return
	}
	return
}
