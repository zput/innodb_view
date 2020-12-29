package view

import (
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
	"gopkg.in/yaml.v2"
)

type IndexPage struct{
	FileAllPage `yaml:"FileAllPage"`

	NDirSlots uint16  `yaml:"NDirSlots"`
	HeapTop uint16 `yaml:"HeapTop"`
	NHeap uint16 `yaml:"NHeap"`
	Free uint16 `yaml:"Free"`
	Garbage uint16 `yaml:"Garbage"`
	LastInsert uint16 `yaml:"LastInsert"`
	Direction uint16 `yaml:"Direction"`
	NDirection uint16 `yaml:"NDirection"`
	NRecs uint16 `yaml:"NRecs"`
	MaxTrxID uint64 `yaml:"MaxTrxID"`
	Level uint16 `yaml:"Level"`
	IndexID uint64 `yaml:"IndexID"`

	LeafNode *TreeNode  `yaml:"LeafNode"`
	NoLeafNode *TreeNode `yaml:"NoLeafNode"`

	RecordSlice []*Record `yaml:"RecordSlice"`
	PageDirectorySlice []PageDirectoryElement `yaml:"PageDirectorySlice"`
}

type Record struct{
	// Variable field lengths(1-2 bytes per var.field) //不定长
	// Nullable field bitmap (1bit per nullable field) //不定长

	// ------------- 5 byte always ------------------------
	InfoFlags InfoFlagsT `yaml:"InfoFlags"` // 4 bits
	NOwned uint8 `yaml:"NOwned"`// 4 bits
	HeapNoIsOrder uint16 `yaml:"HeapNoIsOrder"`// 13 bits
	RecordType uint16 `yaml:"RecordType"`// 3 bits
	NextRecordOffsetRelative int16 `yaml:"NextRecordOffsetRelative"`// 2 byte
}

type InfoFlagsT struct {
	// total (4 bits)
	// saved flag // 1 bit
	// saved flag // 1 bit
	DelFlag uint8 `yaml:"DelFlag"`// 1 bit
	MinFlag uint8 `yaml:"MinFlag"`// 1 bit
}

type PageDirectoryElement struct{
	DirectorySlot uint16 `yaml:"DirectorySlot"` // 2 byte

	NOwned uint8 `yaml:"NOwned"`// 此处物理没有对应的。它实际在record里面
}

func (ip *IndexPage) GetFileType()mysql_define.T_FIL_PAGE_TYPE{
	return mysql_define.T_FIL_PAGE_TYPE(ip.FileAllPage.pageType)
}

func (ip *IndexPage) printPageType() error {
	//prettyFormat, err := json.MarshalIndent(ip, "", "    ")
	prettyFormat, err := yaml.Marshal(ip)
	if err != nil{
		return err
	}
	fmt.Printf("%s\n", string(prettyFormat))
	return nil
}

func (ip *IndexPage) PrintPageType() error {
	ip.printPageType()

	if err := ip.FileAllPage.PrintPageType(); err != nil{
		log.Error(err)
		return err
	}


	return nil
}

func (ip *IndexPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {
	if err := ip.FileAllPage.PageParseFILHeader(buffer); err != nil{
		return err
	}

	return nil
}

func (ip *IndexPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	if err := ip.FileAllPage.PageParseFILTailer(buffer, pageSize); err != nil{
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

	if ip.NoLeafNode, err = getTreeNode(buffer); err != nil{
		log.Error(err)
		return err
	}

	buffer.ExploreBreak()

	if err = ip.parsePageDirectorySlot(buffer, pageSize); err != nil{
		log.Error(err)
		return err
	}

	if err = ip.parseRecords(buffer, pageSize); err != nil{
		log.Error(err)
		return err
	}

	return nil
}

func (ip *IndexPage) parsePageDirectorySlot(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	var isUsingExplore = true
	var err error

	buffer.ExploreBegin()
	if err := buffer.ExploreRetrieve(int(pageSize)-mysql_define.PAGE_DIR-int(ip.NDirSlots)*mysql_define.PAGE_DIR_SLOT_SIZE); err != nil {
		log.Error(err)
		return err
	}

	for i:=int(ip.NDirSlots); i>0; i--{
		var pageDirectoryElement PageDirectoryElement
		pageDirectoryElement.DirectorySlot = buffer.PeekUint16(isUsingExplore)
		if err = buffer.ExploreRetrieve(mysql_define.PAGE_DIR_SLOT_SIZE); err != nil {
			log.Error(err)
			return err
		}
		ip.PageDirectorySlice = append(ip.PageDirectorySlice, pageDirectoryElement)
	}

	buffer.ExploreBreak()

	for index := range ip.PageDirectorySlice{
		buffer.ExploreBegin()
		if err := buffer.ExploreRetrieve(int(ip.PageDirectorySlice[index].DirectorySlot)-mysql_define.REC_N_NEW_EXTRA_BYTES); err != nil {
			log.Error(err)
			return err
		}
		ip.PageDirectorySlice[index].NOwned = buffer.PeekUint8(isUsingExplore)&0x0F
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
	if err := buffer.ExploreRetrieve(int(offset-mysql_define.REC_N_NEW_EXTRA_BYTES)); err != nil {
		log.Error(err)
		return err
	}
	if recordPtr, err = getRecord(buffer, isUsingExplore); err != nil{
		log.Error(err)
		return err
	}
	ip.RecordSlice = append(ip.RecordSlice, recordPtr)
	buffer.ExploreBreak()
	offset += recordPtr.NextRecordOffsetRelative

	for {
		buffer.ExploreBegin()
		if err := buffer.ExploreRetrieve(int(offset-mysql_define.REC_N_NEW_EXTRA_BYTES)); err != nil {
			log.Error(err)
			return err
		}
		if recordPtr, err = getRecord(buffer, isUsingExplore); err != nil{
			log.Error(err)
			return err
		}
		ip.RecordSlice = append(ip.RecordSlice, recordPtr)
		buffer.ExploreBreak()

		if recordPtr.NextRecordOffsetRelative == 0{
			log.Debugf("IndexPage.parseRecords; all record have found in this page")
			break
		}
		offset += recordPtr.NextRecordOffsetRelative
	}
	return nil
}
