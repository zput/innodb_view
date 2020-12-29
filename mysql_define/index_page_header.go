package mysql_define

const PAGE_HEADER = FIL_PAGE_DATA /* index page header starts at this
   offset(除了页面的fil_header) */


const INDEX_PAGE_HEADER_SIZE = 36
const FSEG_HEADER_TOTAL_SIZE = 2*FSEG_HEADER_SIZE

const INDEX_PAGE_BEFORE_RECORD = FIL_PAGE_DATA + INDEX_PAGE_HEADER_SIZE + FSEG_HEADER_TOTAL_SIZE

/** 3.index page header 38-94 **/

/*--------------Index Header---------------*/
const PAGE_N_DIR_SLOTS = 0 /* number of slots in page directory */
const PAGE_HEAP_TOP = 2    /* pointer to record heap top */
const PAGE_N_HEAP = 4      /* number of records in the heap,
   bit 15=flag: new-style compact page format */
const PAGE_FREE = 6         /* pointer to start of page free record list */
const PAGE_GARBAGE = 8      /* number of bytes in deleted records */
const PAGE_LAST_INSERT = 10 /* pointer to the last inserted record, or
   NULL if this info has been reset by a delete,
   for example */
const PAGE_DIRECTION = 12   /* last insert direction: PAGE_LEFT, ... */
const PAGE_N_DIRECTION = 14 /* number of consecutive inserts to the same
   direction */
const PAGE_N_RECS = 16     /* number of user records on the page */
const PAGE_MAX_TRX_ID = 18 /* highest id of a trx which may have modified
   a record on the page; trx_id_t; defined only
   in secondary indexes and in the insert buffer
   tree */
const PAGE_HEADER_PRIV_END = 26 /* end of private data structure of the page
   header which are set in a page create */
/*----*/
const PAGE_LEVEL = 26 /* level of the node in an index tree; the
   leaf level is the level 0.  This field should
   not be written to after page creation. */
const PAGE_INDEX_ID = 28 /* index id where the page belongs.
   This field should not be written to after
   page creation. */
/*--------------Index Header---------------*/

/*-------------- FSEG Header ---------------*/
const PAGE_BTR_SEG_LEAF = 36 /* file segment header for the leaf pages in
   a B-tree: defined only on the root page of a
   B-tree, but not in the root of an ibuf tree */

const PAGE_BTR_SEG_TOP = (36 + FSEG_HEADER_SIZE) /* file segment header for the non-leaf pages
 * in a B-tree: defined only on the root page of
 * a B-tree, but not in the root of an ibuf
 * tree */

//The file segment header points to the inode describing the file segment.
const FSEG_HDR_SPACE = 0   /*!< space id of the inode */
const FSEG_HDR_PAGE_NO = 4 /*!< page number of the inode */
const FSEG_HDR_OFFSET = 8  /*!< byte offset of the inode */

const FSEG_HEADER_SIZE = 10 /*!< Length of the file system
  header, in bytes */


// ??? TODO
const PAGE_BTR_IBUF_FREE_LIST = PAGE_BTR_SEG_LEAF
const PAGE_BTR_IBUF_FREE_LIST_NODE = PAGE_BTR_SEG_LEAF

/* in the place of PAGE_BTR_SEG_LEAF and _TOP
 * there is a free list base node if the page is
 * the root page of an ibuf tree, and at the same
 * place is the free list node if the page is in
 * a free list */
/*-------------- FSEG Header ---------------*/
