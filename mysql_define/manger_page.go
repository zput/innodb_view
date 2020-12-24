package mysql_define

// 下面那个大小112 byte

/* list of pages containing segment
 * headers, where not all the segment
 * header slots are reserved */
/* File space header size */
const   FSP_HEADER_SIZE =          (32 + 5 * FLST_BASE_NODE_SIZE)

/*
 *  Here is table space header/file space header every tablespace have only one at page 0
 *  in Fsp0fsp.h 38-150
 *     */
/*--------------FSP Header(zero-filled for XDES pages) ---> (112)---------------*/

const   FSP_SPACE_ID =             0       /* space id */
const   FSP_NOT_USED =             4       /* this field contained a value up to
   which we know that the modifications
   in the database have been flushed to
   the file space; not used now */
const   FSP_SIZE =                 8       /* Current size of the space in
   pages */
const   FSP_FREE_LIMIT =           12      /* Minimum page number for which the
   free list has not been initialized:
   the pages >= this limit are, by
   definition, free; note that in a
   single-table tablespace where size
   < 64 pages, this number is 64, i.e.,
   we have initialized the space
   about the first extent, but have not
   physically allocated those pages to the
   file */
const   FSP_SPACE_FLAGS =          16      /* fsp_space_t.flags, similar to
   dict_table_t::flags */
const   FSP_FRAG_N_USED =          20      /* number of used pages in the
   FSP_FREE_FRAG list */
const   FSP_FREE =                 24      /* list of free extents */

const   FSP_FREE_FRAG =            (24 + FLST_BASE_NODE_SIZE)
/* list of partially free extents not
 *                                         belonging to any segment */
const   FSP_FULL_FRAG =            (24 + 2 * FLST_BASE_NODE_SIZE)
/* list of full extents not belonging
 *                                         to any segment */
const   FSP_SEG_ID =               (24 + 3 * FLST_BASE_NODE_SIZE)
/* 8 bytes which give the first unused
 *                                         segment id */
const   FSP_SEG_INODES_FULL =      (32 + 3 * FLST_BASE_NODE_SIZE)
/* list of pages containing segment
 *                                         headers, where all the segment inode
 *                                                                                 slots are reserved */
const   FSP_SEG_INODES_FREE =      (32 + 4 * FLST_BASE_NODE_SIZE)
/*--------------FSP Header(zero-filled for XDES pages) ---> (112)---------------*/





















// TODO ???
const   FSP_FREE_ADD =             4       /* this many free extents are added
   to the free list from above
   FSP_FREE_LIMIT at a time */


