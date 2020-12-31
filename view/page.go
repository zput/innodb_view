package view

import (
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/innodb_view/print"
	"github.com/zput/ringbuffer"
)

type ListBaseNode struct {
	FlstLen uint32           `yaml:"FlstLen" self:"FlstLen"`
	First   *PageNoANDOffset `yaml:"First" self:"First"`
	Last    *PageNoANDOffset `yaml:"Last" self:"Last"`
}

type ListNode struct {
	First *PageNoANDOffset `yaml:"First" self:"First"`
	Last  *PageNoANDOffset `yaml:"Last" self:"Last"`
}

type PageNoANDOffset struct {
	PageNo uint32 `yaml:"PageNo" self:"PageNo"`
	Offset uint16 `yaml:"Offset" self:"Offset"`
}

type TreeNode struct {
	SpaceID      uint32           `yaml:"SpaceID" self:"SpaceID"`
	NodePosition *PageNoANDOffset `yaml:"NodePosition" self:"NodePosition"`
}

//---------------------------------------------------------------------//

type IPageParse interface {
	PageParseFILHeader(buffer *ringbuffer.RingBuffer) error
	PageParseFILTailer(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error
	PageParseBody(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error
	PrintPageType() error

	GetFileType() mysql_define.T_FIL_PAGE_TYPE
}

type PageParseFactory struct{}

func (f *PageParseFactory) Create(PageType mysql_define.T_FIL_PAGE_TYPE) IPageParse {
	switch PageType {
	case mysql_define.FIL_COMMON_HEADER_TAILER:
		return new(FileAllPage)
	case mysql_define.FIL_PAGE_TYPE_FSP_HDR:
		return new(FspHeaderPage)
	case mysql_define.FIL_PAGE_INODE:
		return new(INodePage)
	case mysql_define.FIL_PAGE_INDEX:
		return new(IndexPage)
	case mysql_define.FIL_PAGE_TYPE_XDES:

	case mysql_define.FIL_PAGE_TYPE_ALLOCATED:

	}
	return nil
}

// ----------------- Fil header trailer ------------------------------------//

type FileAllPage struct {
	FileHeader  `yaml:"FileHeader" self:""`
	FileTrailer `yaml:"FileTrailer" self:""`
}

type FileHeader struct {
	//FIL_PAGE_SPACE_OR_CHKSUM = 0
	//FIL_PAGE_OFFSET = 4
	//FIL_PAGE_PREV = 8
	//FIL_PAGE_NEXT = 12
	//FIL_PAGE_LSN = 16
	//FIL_PAGE_TYPE = 24
	//FIL_PAGE_FILE_FLUSH_LSN = 26
	//FIL_NULL = 0xFFFFFFFF /*no PAGE_NEXT or PAGE_PREV */
	//FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID = 34

	Checksum                    uint32 `yaml:"Checksum" self:"Checksum"`
	Offset                      uint32 `yaml:"Offset" self:"Offset"`
	PreviousPage                uint32 `yaml:"PreviousPage" self:"PreviousPage"`
	NextPage                    uint32 `yaml:"NextPage" self:"NextPage"`
	LsnForLastPageModeification uint64 `yaml:"LsnForLastPageModeification" self:"LsnForLastPageModeification"`
	PageType                    uint16 `yaml:"PageType" self:"PageType"`
	FlushLSN                    uint64 `yaml:"FlushLSN" self:"FlushLSN"` //(0 except space0 page0) `yaml:"name" self:"name"`
	SpaceID                     uint32 `yaml:"SpaceID" self:"SpaceID"`
}

type FileTrailer struct {
	// FIL_PAGE_END_LSN_OLD_CHKSUM 8 byte
	OldStyleChecksum uint32 `yaml:"OldStyleChecksum" self:"OldStyleChecksum"`
	Lsn              uint32 `yaml:"Lsn" self:"Lsn"`
}

func (fap *FileAllPage) GetFileType() mysql_define.T_FIL_PAGE_TYPE {
	return mysql_define.T_FIL_PAGE_TYPE(fap.PageType)
}

func (fap *FileAllPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {

	var isUsingExplore = true

	buffer.ExploreBegin()

	// TODO optimize this handler error
	// FIL_PAGE_OFFSET
	fap.Checksum = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_OFFSET); err != nil {
		log.Error(err)
		return err
	}

	fap.Offset = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_PREV - mysql_define.FIL_PAGE_OFFSET); err != nil {
		log.Error(err)
		return err
	}

	fap.PreviousPage = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_NEXT - mysql_define.FIL_PAGE_PREV); err != nil {
		log.Error(err)
		return err
	}

	fap.NextPage = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_LSN - mysql_define.FIL_PAGE_NEXT); err != nil {
		log.Error(err)
		return err
	}

	fap.LsnForLastPageModeification = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_TYPE - mysql_define.FIL_PAGE_LSN); err != nil {
		log.Error(err)
		return err
	}

	fap.PageType = buffer.PeekUint16(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_FILE_FLUSH_LSN - mysql_define.FIL_PAGE_TYPE); err != nil {
		log.Error(err)
		return err
	}

	fap.FlushLSN = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID - mysql_define.FIL_PAGE_FILE_FLUSH_LSN); err != nil {
		log.Error(err)
		return err
	}

	fap.SpaceID = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_DATA - mysql_define.FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID); err != nil {
		log.Error(err)
		return err
	}

	buffer.ExploreBreak()

	return nil
}

func (fap *FileAllPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {

	var isUsingExplore = true

	buffer.ExploreBegin()

	if err := buffer.ExploreRetrieve(int(pageSize) - mysql_define.FIL_PAGE_END_LSN_OLD_CHKSUM); err != nil {
		log.Error(err)
		return err
	}
	fap.OldStyleChecksum = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_END_LSN_OLD_CHKSUM / 2); err != nil {
		log.Error(err)
		return err
	}

	fap.Lsn = buffer.PeekUint32(isUsingExplore)

	buffer.ExploreBreak()

	return nil
}

func (fap *FileAllPage) PageParseBody(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	return nil
}

// ------------------------------------------------------- //

func (fap *FileAllPage) printPageType() error {
	log.Debugf("page type value:%d", fap.PageType)

	pageNumber := fap.Offset

	switch mysql_define.T_FIL_PAGE_TYPE(fap.PageType) {

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

func (fap *FileAllPage) generateHumanFormatHeader() []print.PrintFormatT {
	var waitPrintT []print.PrintFormatT
	var currentPosition = mysql_define.FIL_PAGE_SPACE_OR_CHKSUM
	waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "FILE HEADER"))
	waitPrintT = append(waitPrintT, print.Translate(&currentPosition, fap.FileHeader)...)

	return waitPrintT
}

func (fap *FileAllPage) generateHumanFormatTrailer() []print.PrintFormatT {
	var waitPrintT []print.PrintFormatT
	var currentPosition = "N"
	waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "FILE TRAILER"))
	waitPrintT = append(waitPrintT, print.Translate(&currentPosition, fap.FileTrailer)...)

	return waitPrintT
}

func (fap *FileAllPage) PrintPageType() error {

	if err := fap.printPageType(); err != nil {
		return err
	}
	return nil
}

// ----------------- common function ------------------------------------//

func getListBaseNode(buffer *ringbuffer.RingBuffer) (*ListBaseNode, error) {

	var isUsingExplore = true
	var listBaseNode ListBaseNode
	var err error

	listBaseNode.FlstLen = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.PRV_PAGE_NODE); err != nil {
		log.Error(err)
		return nil, err
	}

	listBaseNode.First, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	listBaseNode.Last, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &listBaseNode, nil
}

func getTreeNode(buffer *ringbuffer.RingBuffer) (*TreeNode, error) {

	var isUsingExplore = true
	var treeNode TreeNode
	var err error

	treeNode.SpaceID = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSEG_HDR_PAGE_NO - mysql_define.FSEG_HDR_SPACE); err != nil {
		log.Error(err)
		return nil, err
	}

	treeNode.NodePosition, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &treeNode, nil
}

func getPageNoANDOffset(buffer *ringbuffer.RingBuffer, isUsingExplore bool) (*PageNoANDOffset, error) {

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

func getINodeEntry(buffer *ringbuffer.RingBuffer, isUsingExplore bool) (INodeEntry, error) {

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

	if iNodeEntry.FSegFree, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}
	if iNodeEntry.FSegNotFull, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}
	if iNodeEntry.FSegFull, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}

	iNodeEntry.FSegMagicN = buffer.PeekUint32(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSEG_MAGIC_N); err != nil {
		log.Error(err)
		return iNodeEntry, err
	}

	for i := 0; i < 32; i++ {
		var fSegFragTemp uint32

		if fSegFragTemp, err = getFSegFragArr(buffer, isUsingExplore); err != nil {
			log.Errorf("index[%d]; error[%v]", i, err)
			return iNodeEntry, err
		}
		if fSegFragTemp == 0xffffffff {
			log.Debug("all FSEG_FRAG_ARR_I object have showed in this segment object on INode page")
			if err = buffer.ExploreRetrieve(mysql_define.FSEG_FRAG_ARR_I * (31 - i)); err != nil {
				log.Error(err)
				return iNodeEntry, err
			}
			break
		}

		iNodeEntry.FSegFragSlice = append(iNodeEntry.FSegFragSlice, fSegFragTemp)
	}

	return iNodeEntry, nil
}

func getFSegFragArr(buffer *ringbuffer.RingBuffer, isUsingExplore bool) (fSegFrag uint32, err error) {
	fSegFrag = buffer.PeekUint32(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.FSEG_FRAG_ARR_I); err != nil {
		log.Error(err)
		return
	}
	return
}

func getRecord(buffer *ringbuffer.RingBuffer, isUsingExplore bool) (recordPtr *Record, err error) {
	recordPtr = new(Record)

	InfoFlagsPlusNOwned := buffer.PeekUint8(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.InfoFlagsPlusNOwned); err != nil {
		log.Error(err)
		return
	}
	recordPtr.InfoFlags.DelFlag = (InfoFlagsPlusNOwned & 0x20) >> 5
	recordPtr.InfoFlags.MinFlag = (InfoFlagsPlusNOwned & 0x10) >> 4
	recordPtr.NOwned = InfoFlagsPlusNOwned & 0x0F

	HeapNoIsOrderPlusRecordType := buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.HeapNoPlusRecordType); err != nil {
		log.Error(err)
		return
	}
	recordPtr.HeapNoIsOrder = (HeapNoIsOrderPlusRecordType & 0xFFF8) >> 3
	recordPtr.RecordType = HeapNoIsOrderPlusRecordType & 0x7

	nextRecord := buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.NextRecord); err != nil {
		log.Error(err)
		return
	}
	recordPtr.NextRecordOffsetRelative = int16(nextRecord)
	return
}
