package mysql_define



/** 1.file page header 1-38 **/
/** The byte offsets on a file page for various variables @{
 *  */
const   FIL_PAGE_SPACE_OR_CHKSUM =  0	/*!< in < MySQL-4.0.14 space id the
  page belongs to (== 0) but in later
  versions the 'new' checksum of the
  page */
const   FIL_PAGE_OFFSET = 		4	/*!< page offset inside space */
const   FIL_PAGE_PREV = 		8	/*!< if there is a 'natural'
  predecessor of the page, its
  offset.  Otherwise FIL_NULL.
  This field is not set on BLOB
  pages, which are stored as a
  singly-linked list.  See also
  FIL_PAGE_NEXT. */
const   FIL_PAGE_NEXT = 		12	/*!< if there is a 'natural' successor
  of the page, its offset.
  Otherwise FIL_NULL.
  B-tree index pages
  (FIL_PAGE_TYPE contains FIL_PAGE_INDEX)
  on the same PAGE_LEVEL are maintained
  as a doubly linked list via
  FIL_PAGE_PREV and FIL_PAGE_NEXT
  in the collation order of the
  smallest user record on each page. */
const   FIL_PAGE_LSN = 		16	/*!< lsn of the end of the newest
  modification log record to the page */

// ----------------this value about file page type is below.------------------
const 	 FIL_PAGE_TYPE = 		24	/*!< file page type: FIL_PAGE_INDEX,...,
  2 bytes.

  The contents of this field can only
  be trusted in the following case:
  if the page is an uncompressed
  B-tree index page, then it is
  guaranteed that the value is
  FIL_PAGE_INDEX.
  The opposite does not hold.

  In tablespaces created by
  MySQL/InnoDB 5.1.7 or later, the
  contents of this field is valid
  for all uncompressed pages. */
const   FIL_PAGE_FILE_FLUSH_LSN = 	26	/*!< this is only defined for the
  first page of the system tablespace:
  the file has been flushed to disk
  at least up to this LSN. For
  FIL_PAGE_COMPRESSED pages, we store
  the compressed page control information
  in these 8 bytes. */
const 	 FIL_NULL =  0xFFFFFFFF       /*no PAGE_NEXT or PAGE_PREV */
/** starting from 4.1.x this contains the space id of the page */
const   FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID =   34

const   FIL_PAGE_SPACE_ID =   FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID

const   FIL_PAGE_DATA = 		38	/*!< start of the data on the page */




/** 2.page type value **/

/** File page types (values of FIL_PAGE_TYPE) @{
 *  */
const   FIL_PAGE_INDEX = 		17855	/*!< B-tree node noraml data page*/
const   FIL_PAGE_RTREE = 		17854	/*!< R-tree node */
const   FIL_PAGE_UNDO_LOG = 	2	/*!< Undo log page */
const   FIL_PAGE_INODE = 		3	/*!< Index node */
const   FIL_PAGE_IBUF_FREE_LIST = 	4	/*!< Insert buffer free list */
/* File page types introduced in MySQL/InnoDB 5.1.7 */
const   FIL_PAGE_TYPE_ALLOCATED = 	0	/*!< Freshly allocated page */
const   FIL_PAGE_IBUF_BITMAP = 	5	/*!< Insert buffer bitmap */
const   FIL_PAGE_TYPE_SYS = 	6	/*!< System page */
const   FIL_PAGE_TYPE_TRX_SYS = 	7	/*!< Transaction system data */
const   FIL_PAGE_TYPE_FSP_HDR = 	8	/*!< File space header */
const   FIL_PAGE_TYPE_XDES = 	9	/*!< Extent descriptor page */
const   FIL_PAGE_TYPE_BLOB = 	10	/*!< Uncompressed BLOB page */
const   FIL_PAGE_TYPE_ZBLOB = 	11	/*!< First compressed BLOB page */
const   FIL_PAGE_TYPE_ZBLOB2 = 	12	/*!< Subsequent compressed BLOB page */
const   FIL_PAGE_TYPE_UNKNOWN = 	13	/*!< In old tablespaces, garbage
  in FIL_PAGE_TYPE is replaced with this
  value when flushing pages. */
const   FIL_PAGE_COMPRESSED = 	14	/*!< Compressed page */
const   FIL_PAGE_ENCRYPTED = 	15	/*!< Encrypted page */
const   FIL_PAGE_COMPRESSED_AND_ENCRYPTED =  16
/*!< Compressed and Encrypted page */
const   FIL_PAGE_ENCRYPTED_RTREE =  17	/*!< Encrypted R-tree page */



//...
//...
//...



/** 5.page end/trailer **/
const   FIL_PAGE_END_LSN_OLD_CHKSUM =  8	/*!< the low 4 bytes of this are used
  to store the page checksum, the
  last 4 bytes should be identical
  to the last 4 bytes of FIL_PAGE_LSN */



const   FIL_PAGE_DATA_END = 	8	/*!< size of the page trailer */
