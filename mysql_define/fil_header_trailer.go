package mysql_define

type PAGE_SIZE int

const (
	_                        = iota
	PAGE_SIZE_16_K PAGE_SIZE = 16 * 1024 * iota
	PAGE_SIZE_32_K
	PAGE_SIZE_48_K
	PAGE_SIZE_64_K
)

const FIL_NULL = 0xFFFFFFFF /*no PAGE_NEXT or PAGE_PREV */

const FIL_PAGE_DATA = 38 /*!< start of the data on the page */

const FIL_PAGE_TRAILER_BEFORE_SIZE = 16376

/** 1.file page header 1-38 **/
/** The byte offsets on a file page for various variables @{
 *  */
const FIL_PAGE_SPACE_OR_CHKSUM = 0 /*!< in < MySQL-4.0.14 space id the
  page belongs to (== 0) but in later
  versions the 'new' checksum of the
  page */
const FIL_PAGE_OFFSET = 4 /*!< page offset inside space */
const FIL_PAGE_PREV = 8   /*!< if there is a 'natural'
  predecessor of the page, its
  offset.  Otherwise FIL_NULL.
  This field is not set on BLOB
  pages, which are stored as a
  singly-linked list.  See also
  FIL_PAGE_NEXT. */
const FIL_PAGE_NEXT = 12 /*!< if there is a 'natural' successor
  of the page, its offset.
  Otherwise FIL_NULL.
  B-tree index pages
  (FIL_PAGE_TYPE contains FIL_PAGE_INDEX)
  on the same PAGE_LEVEL are maintained
  as a doubly linked list via
  FIL_PAGE_PREV and FIL_PAGE_NEXT
  in the collation order of the
  smallest user record on each page. */
const FIL_PAGE_LSN = 16 /*!< lsn of the end of the newest
  modification log record to the page */

// ----------------this value about file page type is below.------------------
const FIL_PAGE_TYPE = 24 /*!< file page type: FIL_PAGE_INDEX,...,
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
const FIL_PAGE_FILE_FLUSH_LSN = 26 /*!< this is only defined for the
  first page of the system tablespace:
  the file has been flushed to disk
  at least up to this LSN. For
  FIL_PAGE_COMPRESSED pages, we store
  the compressed page control information
  in these 8 bytes. */

/** starting from 4.1.x this contains the space id of the page */
const FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID = 34

const FIL_PAGE_SPACE_ID = FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID

/** 2.page type value **/

type T_FIL_PAGE_TYPE int

const (
	FIL_COMMON_HEADER_TAILER T_FIL_PAGE_TYPE = -1

	/** File page types (values of FIL_PAGE_TYPE) @{
	 *  */
	FIL_PAGE_INDEX          T_FIL_PAGE_TYPE = 17855 /*!< B-tree node noraml data page*/
	FIL_PAGE_RTREE          T_FIL_PAGE_TYPE = 17854 /*!< R-tree node */
	FIL_PAGE_UNDO_LOG       T_FIL_PAGE_TYPE = 2     /*!< Undo log page */
	FIL_PAGE_INODE          T_FIL_PAGE_TYPE = 3     /*!< Index node */
	FIL_PAGE_IBUF_FREE_LIST T_FIL_PAGE_TYPE = 4     /*!< Insert buffer free list */
	/* File page types introduced in MySQL/InnoDB 5.1.7 */
	FIL_PAGE_TYPE_ALLOCATED T_FIL_PAGE_TYPE = 0  /*!< Freshly allocated page */
	FIL_PAGE_IBUF_BITMAP    T_FIL_PAGE_TYPE = 5  /*!< Insert buffer bitmap */
	FIL_PAGE_TYPE_SYS       T_FIL_PAGE_TYPE = 6  /*!< System page */
	FIL_PAGE_TYPE_TRX_SYS   T_FIL_PAGE_TYPE = 7  /*!< Transaction system data */
	FIL_PAGE_TYPE_FSP_HDR   T_FIL_PAGE_TYPE = 8  /*!< File space header */
	FIL_PAGE_TYPE_XDES      T_FIL_PAGE_TYPE = 9  /*!< Extent descriptor page */
	FIL_PAGE_TYPE_BLOB      T_FIL_PAGE_TYPE = 10 /*!< Uncompressed BLOB page */
	FIL_PAGE_TYPE_ZBLOB     T_FIL_PAGE_TYPE = 11 /*!< First compressed BLOB page */
	FIL_PAGE_TYPE_ZBLOB2    T_FIL_PAGE_TYPE = 12 /*!< Subsequent compressed BLOB page */
	FIL_PAGE_TYPE_UNKNOWN   T_FIL_PAGE_TYPE = 13 /*!< In old tablespaces, garbage
	  in FIL_PAGE_TYPE is replaced with this
	  value when flushing pages. */
	FIL_PAGE_COMPRESSED               T_FIL_PAGE_TYPE = 14 /*!< Compressed page */
	FIL_PAGE_ENCRYPTED                T_FIL_PAGE_TYPE = 15 /*!< Encrypted page */
	FIL_PAGE_COMPRESSED_AND_ENCRYPTED T_FIL_PAGE_TYPE = 16

	/*!< Compressed and Encrypted page */
	FIL_PAGE_ENCRYPTED_RTREE T_FIL_PAGE_TYPE = 17 /*!< Encrypted R-tree page */
)

var statusText = map[T_FIL_PAGE_TYPE]string{
	FIL_PAGE_INDEX:          "INDEX page",
	FIL_PAGE_INODE:          "INODE page(segment object)",
	FIL_PAGE_TYPE_ALLOCATED: "Freshly allocated page",
	FIL_PAGE_TYPE_FSP_HDR:   "File space header",
	FIL_PAGE_TYPE_XDES:      "XDES page",
	FIL_PAGE_IBUF_BITMAP:    "INSERT buffer bitmap page",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code T_FIL_PAGE_TYPE) string {
	return statusText[code]
}

//...
//...
//...

/** 5.page end/trailer **/
const FIL_PAGE_END_LSN_OLD_CHKSUM = 8 /*!< the low 4 bytes of this are used
  to store the page checksum, the
  last 4 bytes should be identical
  to the last 4 bytes of FIL_PAGE_LSN */

const FIL_PAGE_DATA_END = 8 /*!< size of the page trailer */
