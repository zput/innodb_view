package view

import (
	"encoding/json"
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
)

type INodePage struct{
	FileAllPage
	List ListNode
	INodeEntrySlice [85]INodeEntry
}

type INodeEntry struct{
	FSegID uint64
	FSegNotFullNUsed uint64
	FSegFree *ListBaseNode
	FSegNotFull *ListBaseNode
	FSegFull *ListBaseNode
	FSegMagicN uint32
	FSegFragSlice [32]uint32
}

func (inp *INodePage) GetFileType()mysql_define.T_FIL_PAGE_TYPE{
	return mysql_define.T_FIL_PAGE_TYPE(inp.FileAllPage.pageType)
}

func (inp *INodePage) printPageType() error {
	prettyFormat, err := json.MarshalIndent(inp, "", "    ")
	if err != nil{
		return err
	}
	fmt.Printf("%s", string(prettyFormat))
	return nil
}

func (inp *INodePage) PrintPageType() error {
	inp.printPageType()

	if err := inp.FileAllPage.PrintPageType(); err != nil{
		log.Error(err)
		return err
	}


	return nil
}

func (inp *INodePage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {
	if err := inp.FileAllPage.PageParseFILHeader(buffer); err != nil{
		return err
	}

	return nil
}

func (inp *INodePage) PageParseFILTailer(buffer *ringbuffer.RingBuffer) error {
	if err := inp.FileAllPage.PageParseFILTailer(buffer); err != nil{
		return err
	}

	return nil
}

func (inp *INodePage) PageParseBody(buffer *ringbuffer.RingBuffer) error {

	var isUsingExplore = true
	var err error

	buffer.ExploreBegin()

	if err := buffer.ExploreRetrieve(mysql_define.FIL_PAGE_DATA); err != nil {
		log.Error(err)
		return err
	}

	inp.List.First, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil{
		log.Error(err)
		return err
	}

	inp.List.Last, err = getPageNoANDOffset(buffer, isUsingExplore)
	if err != nil{
		log.Error(err)
		return err
	}

	for i:=0; i<85;i++{
		if inp.INodeEntrySlice[i], err = getINodeEntry(buffer, isUsingExplore); err != nil{
			log.Errorf("index[%d]; error[%v]", i, err)
			return err
		}
	}

	buffer.ExploreBreak()

	return nil
}
