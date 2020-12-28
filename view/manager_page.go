package view

import (
	"encoding/json"
	"fmt"
	"github.com/zput/innodb_view/log"
	"github.com/zput/innodb_view/mysql_define"
	"github.com/zput/ringbuffer"
)

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


