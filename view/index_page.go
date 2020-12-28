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

func (ip *IndexPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer) error {
	if err := ip.FileAllPage.PageParseFILTailer(buffer); err != nil{
		return err
	}

	return nil
}

func (ip *IndexPage) PageParseBody(buffer *ringbuffer.RingBuffer) error {

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
	return nil
}
