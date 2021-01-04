package view

import (
	"encoding/json"
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/print"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
)

// ----------------- FspHeaderPage ------------------------------------//
type FspHeaderPage struct {
	FileAllPage `yaml:"FileAllPage" self:"FileAllPage"`
	FspHeader `yaml:"FspHeader" self:"FspHeader"`
}

type FspHeader struct {
	FspSpaceID    uint32 `yaml:"FspSpaceID" self:"FspSpaceID"`
	FspNotUsed    uint32 `yaml:"FspNotUsed" self:"FspNotUsed"`
	FspSize       uint32 `yaml:"FspSize" self:"FspSize"`
	FspFreeLimit  uint32 `yaml:"FspFreeLimit" self:"FspFreeLimit"`
	FspSpaceFlags uint32 `yaml:"FspSpaceFlags" self:"FspSpaceFlags"`
	FspFragNUsed  uint32 `yaml:"FspFragNUsed" self:"FspFragNUsed"`

	FspFree     *ListBaseNode `yaml:"FspFree" self:"FspFree"`
	FspFreeFrag *ListBaseNode `yaml:"FspFreeFrag" self:"FspFreeFrag"`
	FspFullFrag *ListBaseNode `yaml:"FspFullFrag" self:"FspFullFrag"`

	FspSegID uint64 `yaml:"FspSegID" self:"FspSegID"`

	FspSegInodesFull *ListBaseNode `yaml:"FspSegInodesFull" self:"FspSegInodesFull"`
	FspSegInodesFree *ListBaseNode `yaml:"FspSegInodesFree" self:"FspSegInodesFree"`
}

func (fhp *FspHeaderPage) GetFileType() mysql_define.T_FIL_PAGE_TYPE {
	return mysql_define.T_FIL_PAGE_TYPE(fhp.FileAllPage.PageType)
}

func (fhp *FspHeaderPage) PageParseFILHeader(buffer *ringbuffer.RingBuffer) error {
	if err := fhp.FileAllPage.PageParseFILHeader(buffer); err != nil {
		return err
	}

	return nil
}

func (fhp *FspHeaderPage) PageParseFILTailer(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {
	if err := fhp.FileAllPage.PageParseFILTailer(buffer, pageSize); err != nil {
		return err
	}

	return nil
}

func (fhp *FspHeaderPage) PageParseBody(buffer *ringbuffer.RingBuffer, pageSize mysql_define.PAGE_SIZE) error {

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
	if fhp.FspFree, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return err
	}
	if fhp.FspFreeFrag, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return err
	}
	if fhp.FspFullFrag, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return err
	}

	fhp.FspSegID = buffer.PeekUint64(isUsingExplore)
	if err := buffer.ExploreRetrieve(mysql_define.FSP_SEG_INODES_FULL - mysql_define.FSP_SEG_ID); err != nil {
		log.Error(err)
		return err
	}

	if fhp.FspSegInodesFull, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return err
	}
	if fhp.FspSegInodesFree, err = getListBaseNode(buffer); err != nil {
		log.Error(err)
		return err
	}

	buffer.ExploreBreak()

	return nil
}

// --------------- inner method function ----------------- //
func (fhp *FspHeaderPage) printPageType() error {
	prettyFormat, err := json.MarshalIndent(fhp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("%s", string(prettyFormat))
	return nil
}

func (fhp *FspHeaderPage) generateHumanFormat() []print.PrintFormatT {
	var waitPrintT []print.PrintFormatT
	var currentPosition int

	waitPrintT = append(waitPrintT, fhp.FileAllPage.generateHumanFormatHeader()...)

	waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "index page:FSP header"))
	currentPosition = mysql_define.FIL_PAGE_DATA
	currentPosition *= 8
	waitPrintT = append(waitPrintT, print.Translate(&currentPosition, fhp.FspHeader)...)

	// TODO
	//waitPrintT = append(waitPrintT, *print.NewPrintFormatT(print.PrintDivideSignBlock, "index page:entry(0-84)"))
	//currentPosition = mysql_define.FIL_PAGE_DATA+mysql_define.FSEG_INODE_PAGE_NODE
	//currentPosition *= 8
	//waitPrintT = append(waitPrintT, print.Translate(&currentPosition, fhp.INodeEntrySlice)...)

	waitPrintT = append(waitPrintT, fhp.FileAllPage.generateHumanFormatTrailer()...)

	return waitPrintT
}

func (fhp *FspHeaderPage) PrintPageType() error {
	fmt.Printf("%s\n", print.PrintFun(fhp.generateHumanFormat()))

	fmt.Println()

	//fhp.printPageType()

	if err := fhp.FileAllPage.PrintPageType(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
