package mysql_define

const FIL_ADDR_SIZE = 6 /* address size is 6 bytes */

const FLST_BASE_NODE_SIZE = (4 + 2*FIL_ADDR_SIZE) /*
 *                list--length   :4
 * FIL_ADDR_SIZE prv page node   :4
 *                   offset      :2
 * FIL_ADDR_SIZE nxt page node   :4
 * 					 offset      :2
 */

const (
	LIST_LENGHT    = 0
	PRV_PAGE_NODE  = 4
	PRV_OFFSET     = 8
	NEXT_PAGE_NODE = 10
	NEXT_OFFSET    = 14

	LIST_BASE_NODE_SIZE = 16
)
