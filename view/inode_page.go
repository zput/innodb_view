package view

import (
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
	"gopkg.in/yaml.v2"
)

type INodePage struct{
	FileAllPage `yaml:"FileAllPage"`
	List ListNode `yaml:"List"`
	INodeEntrySlice []INodeEntry `yaml:"INodeEntrySlice" json:"INodeEntrySlice,omitempty"`
	//INodeEntrySlice [85]INodeEntry `yaml:"INodeEntrySlice" json:"INodeEntrySlice,omitempty"`
}

type INodeEntry struct{
	FSegID uint64 `yaml:"FSegID"`
	FSegNotFullNUsed uint64 `yaml:"FSegNotFullNUsed"`
	FSegFree *ListBaseNode `yaml:"FSegFree"`
	FSegNotFull *ListBaseNode `yaml:"FSegNotFull"`
	FSegFull *ListBaseNode `yaml:"FSegFull"`
	FSegMagicN uint32 `yaml:"FSegMagicN"`
	FSegFragSlice []uint32 `yaml:"FSegFragSlice" json:"FSegFragSlice,omitempty"`
	//FSegFragSlice [32]uint32 `yaml:"FSegFragSlice"`
}

func (inp *INodePage) GetFileType()mysql_define.T_FIL_PAGE_TYPE{
	return mysql_define.T_FIL_PAGE_TYPE(inp.FileAllPage.pageType)
}

func (inp *INodePage) printPageType() error {
	//prettyFormat, err := json.MarshalIndent(inp, "", "    ")
	prettyFormat, err := yaml.Marshal(inp)
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

func (inp *INodePage) PageParseFILTailer(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	if err := inp.FileAllPage.PageParseFILTailer(buffer, pageSize); err != nil{
		return err
	}

	return nil
}

func (inp *INodePage) PageParseBody(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {

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
		var iNodeEntryTmp INodeEntry
		if iNodeEntryTmp, err = getINodeEntry(buffer, isUsingExplore); err != nil{
			log.Errorf("index[%d]; error[%v]", i, err)
			return err
		}
		if iNodeEntryTmp.FSegID <= 0{
			log.Debug("all segment object have showed in this INode page")
			break
		}

		inp.INodeEntrySlice = append(inp.INodeEntrySlice, iNodeEntryTmp)
	}

	buffer.ExploreBreak()

	return nil
}
