package mysql_define

/** 4.INFIMUM && SUPREMUM **/

/* Number of extra bytes in a new-style record,
 * in addition to the data and the offsets */
const REC_N_NEW_EXTRA_BYTES = 5 //new-style记录扩展字节

const InfoFlagsPlusNOwned = 1  // 1 byte
const HeapNoPlusRecordType = 2 // 2 byte
const NextRecord = 2           // 2 byte

/* Record status values */
const REC_STATUS_ORDINARY = 0 // 普通记录
const REC_STATUS_NODE_PTR = 1 // 非叶子结点带指针
const REC_STATUS_INFIMUM = 2  // Infimum
const REC_STATUS_SUPREMUM = 3 // Supremum

/* The following four constants are needed in page0zip.cc in order to
 * efficiently compress and decompress pages. */

/* The offset of heap_no in a compact record */
const REC_NEW_HEAP_NO = 4

/* The shift of heap_no in a compact record.
 * The status is stored in the low-order bits. */
const REC_HEAP_NO_SHIFT = 3

/* Length of a B-tree node pointer, in bytes */
const REC_NODE_PTR_SIZE = 4

/*----*/
const PAGE_DATA = (PAGE_HEADER + 36 + 2*FSEG_HEADER_SIZE)

/* start of data on the page */
const PAGE_NEW_INFIMUM = (PAGE_DATA + REC_N_NEW_EXTRA_BYTES)

/* offset of the page infimum record on a
 *                                 new-style compact page */
const PAGE_NEW_SUPREMUM = (PAGE_DATA + 2*REC_N_NEW_EXTRA_BYTES + 8)

/* offset of the page supremum record on a
 *                                 new-style compact page */
const PAGE_NEW_SUPREMUM_END = (PAGE_NEW_SUPREMUM + 8)

/* offset of the page supremum record end on
 *                                 a new-style compact page */

//------------------------** 6.PAGE DIRECTORY **/ ------------------------------------

/* Offset of the directory start down from the page end. We call the
 * slot with the highest file address directory start, as it points to
 * the first record in the list of records. */
const PAGE_DIR = FIL_PAGE_DATA_END

/* We define a slot in the page directory as two bytes */
const PAGE_DIR_SLOT_SIZE = 2

/* The offset of the physically lower end of the directory, counted from
 * page end, when the page is empty */
const PAGE_EMPTY_DIR_START = (PAGE_DIR + 2*PAGE_DIR_SLOT_SIZE)

/* The maximum and minimum number of records owned by a directory slot. The
 * number may drop below the minimum in the first and the last slot in the
 * directory. */
const PAGE_DIR_SLOT_MAX_N_OWNED = 8
const PAGE_DIR_SLOT_MIN_N_OWNED = 4

const FAIL_CHK = 8
const FAIL_LSN = 4
