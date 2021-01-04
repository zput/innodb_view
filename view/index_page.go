package view

import (
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/innodb_view/print"
	"github.com/zput/ringbuffer"
	//"gopkg.in/yaml.v2"
	"github.com/naoina/toml"
)

type IndexPage struct {
	FileAllPage `yaml:"FileAllPage" self:"FileAllPage"`

	IndexHeader `yaml:"IndexHeader" self:"IndexHeader"`
	FSegHeader  `yaml:"FSegHeader" self:"FSegHeader"`

	IndexRecord `yaml:"IndexRecord" self:"IndexRecord"`
}

type IndexHeader struct {
	NDirSlots  uint16 `yaml:"NDirSlots" self:"NDirSlots"`
	HeapTop    uint16 `yaml:"HeapTop" self:"HeapTop"`
	NHeap      uint16 `yaml:"NHeap" self:"NHeap"`
	Free       uint16 `yaml:"Free" self:"Free"`
	Garbage    uint16 `yaml:"Garbage" self:"Garbage"`
	LastInsert uint16 `yaml:"LastInsert" self:"LastInsert"`
	Direction  uint16 `yaml:"Direction" self:"Direction"`
	NDirection uint16 `yaml:"NDirection" self:"NDirection"`
	NRecs      uint16 `yaml:"NRecs" self:"NRecs"`
	MaxTrxID   uint64 `yaml:"MaxTrxID" self:"MaxTrxID"`
	Level      uint16 `yaml:"Level" self:"Level"`
	IndexID    uint64 `yaml:"IndexID" self:"IndexID"`
}

type FSegHeader struct {
	LeafNode   *TreeNode `yaml:"LeafNode" self:"LeafNode"`
	NoLeafNode *TreeNode `yaml:"NoLeafNode" self:"NoLeafNode"`
}

type IndexRecord struct {
	RecordSlice        []*Record              `yaml:"RecordSlice" self:"RecordSlice"`
	PageDirectorySlice []PageDirectoryElement `yaml:"PageDirectorySlice" self:"PageDirectorySlice"`
}

type Record struct {
	// Variable field lengths(1-2 bytes per var.field) //不定长
	// Nullable field bitmap (1bit per nullable field) //不定长

	// ------------- 5 byte always ------------------------
	InfoFlags                InfoFlagsT `yaml:"InfoFlags" self:"InfoFlags,4"`                               // 4 bits
	NOwned                   uint8      `yaml:"NOwned" self:"NOwned,4"`                                     // 4 bits
	HeapNoIsOrder            uint16     `yaml:"HeapNoIsOrder" self:"HeapNoIsOrder,13"`                       // 13 bits
	RecordType               uint16     `yaml:"RecordType" self:"RecordType,3"`                             // 3 bits
	NextRecordOffsetRelative int16      `yaml:"NextRecordOffsetRelative" self:"NextRecordOffsetRelative"` // 2 byte
}

type InfoFlagsT struct {
	// total (4 bits)
	// saved flag // 1 bit
	// saved flag // 1 bit
	SaveFlag1 uint8 `yaml:"SaveFlag1" self:"SaveFlag1,1"` // 1 bit
	SaveFlag2 uint8 `yaml:"SaveFlag2" self:"SaveFlag2,1"` // 1 bit
	DelFlag uint8 `yaml:"DelFlag" self:"DelFlag,1"` // 1 bit
	MinFlag uint8 `yaml:"MinFlag" self:"MinFlag,1"` // 1 bit
}

type PageDirectoryElement struct {
	DirectorySlot uint16 `yaml:"DirectorySlot" self:"DirectorySlot"` // 2 byte

	NOwned uint8 `yaml:"NOwned" self:"NOwned,0"` // 此处物理没有对应的。它实际在record里面
}

func (ip *IndexPage) GetFileType() mysql_define.T_FIL_PAGE_TYPE {
	return mysql_define.T_FIL_PAGE_TYPE(ip.FileAllPage.PageType)
}

func (ip *IndexPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {
	if err := ip.FileAllPage.PageParseFILHeader(buffer); err != nil {
		return err
	}

	return nil
}

func (ip *IndexPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	if err := ip.FileAllPage.PageParseFILTailer(buffer, pageSize); err != nil {
		return err
	}

	return nil
}

func (ip *IndexPage) PageParseBody(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {

	var isUsingExplore = true
	var err error

	buffer.ExploreBegin()

	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_DATA); err != nil {
		log.Error(err)
		return err
	}

	ip.NDirSlots = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_HEAP_TOP - mysql_define.PAGE_N_DIR_SLOTS); err != nil {
		log.Error(err)
		return err
	}

	ip.HeapTop = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_N_HEAP - mysql_define.PAGE_HEAP_TOP); err != nil {
		log.Error(err)
		return err
	}

	ip.NHeap = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_FREE - mysql_define.PAGE_N_HEAP); err != nil {
		log.Error(err)
		return err
	}

	ip.Free = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_GARBAGE - mysql_define.PAGE_FREE); err != nil {
		log.Error(err)
		return err
	}

	ip.Garbage = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_LAST_INSERT - mysql_define.PAGE_GARBAGE); err != nil {
		log.Error(err)
		return err
	}

	ip.LastInsert = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_DIRECTION - mysql_define.PAGE_LAST_INSERT); err != nil {
		log.Error(err)
		return err
	}

	ip.Direction = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_N_RECS - mysql_define.PAGE_DIRECTION); err != nil {
		log.Error(err)
		return err
	}

	ip.NRecs = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_MAX_TRX_ID - mysql_define.PAGE_N_RECS); err != nil {
		log.Error(err)
		return err
	}

	ip.MaxTrxID = buffer.PeekUint64(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_LEVEL - mysql_define.PAGE_MAX_TRX_ID); err != nil {
		log.Error(err)
		return err
	}

	ip.Level = buffer.PeekUint16(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_INDEX_ID - mysql_define.PAGE_LEVEL); err != nil {
		log.Error(err)
		return err
	}

	ip.IndexID = buffer.PeekUint64(isUsingExplore)
	if err = buffer.ExploreRetrieve(mysql_define.PAGE_BTR_SEG_LEAF - mysql_define.PAGE_INDEX_ID); err != nil {
		log.Error(err)
		return err
	}

	// ----------------------- FSEG header -------------------//
	if ip.LeafNode, err = getTreeNode(buffer); err != nil {
		log.Error(err)
		return err
	}

	if ip.NoLeafNode, err = getTreeNode(buffer); err != nil {
		log.Error(err)
		return err
	}

	buffer.ExploreBreak()

	if err = ip.parsePageDirectorySlot(buffer, pageSize); err != nil {
		log.Error(err)
		return err
	}

	if err = ip.parseRecords(buffer, pageSize); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// --------------- inner method function ----------------- //
func (ip *IndexPage) parsePageDirectorySlot(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	var isUsingExplore = true
	var err error

	buffer.ExploreBegin()
	if err := buffer.ExploreRetrieve(int(pageSize) - mysql_define.PAGE_DIR - int(ip.NDirSlots)*mysql_define.PAGE_DIR_SLOT_SIZE); err != nil {
		log.Error(err)
		return err
	}

	for i := int(ip.NDirSlots); i > 0; i-- {
		var pageDirectoryElement PageDirectoryElement
		pageDirectoryElement.DirectorySlot = buffer.PeekUint16(isUsingExplore)
		if err = buffer.ExploreRetrieve(mysql_define.PAGE_DIR_SLOT_SIZE); err != nil {
			log.Error(err)
			return err
		}
		ip.PageDirectorySlice = append(ip.PageDirectorySlice, pageDirectoryElement)
	}

	buffer.ExploreBreak()

	for index := range ip.PageDirectorySlice {
		buffer.ExploreBegin()
		if err := buffer.ExploreRetrieve(int(ip.PageDirectorySlice[index].DirectorySlot) - mysql_define.REC_N_NEW_EXTRA_BYTES); err != nil {
			log.Error(err)
			return err
		}
		ip.PageDirectorySlice[index].NOwned = buffer.PeekUint8(isUsingExplore) & 0x0F
		buffer.ExploreBreak()
	}

	return nil
}

func (ip *IndexPage) parseRecords(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	var isUsingExplore = true
	var err error
	var offset int16 = mysql_define.INDEX_PAGE_BEFORE_RECORD + mysql_define.REC_N_NEW_EXTRA_BYTES
	var recordPtr *Record

	// get infimum record
	buffer.ExploreBegin()
	if err := buffer.ExploreRetrieve(int(offset - mysql_define.REC_N_NEW_EXTRA_BYTES)); err != nil {
		log.Error(err)
		return err
	}
	if recordPtr, err = getRecord(buffer, isUsingExplore); err != nil {
		log.Error(err)
		return err
	}
	ip.RecordSlice = append(ip.RecordSlice, recordPtr)
	buffer.ExploreBreak()
	offset += recordPtr.NextRecordOffsetRelative

	for {
		buffer.ExploreBegin()
		if err := buffer.ExploreRetrieve(int(offset - mysql_define.REC_N_NEW_EXTRA_BYTES)); err != nil {
			log.Error(err)
			return err
		}
		if recordPtr, err = getRecord(buffer, isUsingExplore); err != nil {
			log.Error(err)
			return err
		}
		ip.RecordSlice = append(ip.RecordSlice, recordPtr)
		buffer.ExploreBreak()

		if recordPtr.NextRecordOffsetRelative == 0 {
			log.Debugf("IndexPage.parseRecords; all record have found in this page")
			break
		}
		offset += recordPtr.NextRecordOffsetRelative
	}
	return nil
}

func (ip *IndexPage) printPageType() error {
	//prettyFormat, err := json.MarshalIndent(ip, "", "    ")
	//prettyFormat, err := yaml.Marshal(ip)
	prettyFormat, err := toml.Marshal(ip)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(prettyFormat))
	return nil
}

func (ip *IndexPage) generateHumanFormat() []print.PrintFormatT {
	var waitPrintT []print.PrintFormatT
	var currentPosition int

	waitPrintT = append(waitPrintT, ip.FileAllPage.generateHumanFormatHeader()...)

	waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "index page:index header"))
	currentPosition = mysql_define.FIL_PAGE_DATA
	currentPosition *= 8
	waitPrintT = append(waitPrintT, print.Translate(&currentPosition, ip.IndexHeader)...)

	waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "index page:FSEG header"))
	currentPosition = mysql_define.FIL_PAGE_DATA+mysql_define.INDEX_PAGE_HEADER_SIZE
	currentPosition *= 8
	waitPrintT = append(waitPrintT, print.Translate(&currentPosition, ip.FSegHeader)...)

	waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "index page:All records"))
	currentPosition = mysql_define.INDEX_PAGE_BEFORE_RECORD
	currentPosition *= 8
	waitPrintT = append(waitPrintT, print.Translate(&currentPosition, ip.IndexRecord)...)

	waitPrintT = append(waitPrintT, ip.FileAllPage.generateHumanFormatTrailer()...)

	return waitPrintT
}

func (ip *IndexPage) PrintPageType() error {

	fmt.Printf("%s\n", print.PrintFun(ip.generateHumanFormat()))

	fmt.Println()

	ip.printPageType()

	if err := ip.FileAllPage.PrintPageType(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
