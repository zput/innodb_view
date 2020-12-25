package mysql_define

const FSEG_INODE_PAGE_NODE = 12 /* INODE页的链表节点，记录前后Inode Page的位置，
BaseNode记录在头Page的FSP_SEG_INODES_FULL或者FSP_SEG_INODES_FREE字段。*/

const INODE_ENTRY_SIZE = 192 /* this size is 192 for every inode entry(segment object)*/

/*--------------------------- Inode Entry struct ---------------------------------------*/

const FSEG_ID = 8 // 该Inode归属的Segment ID，若值为0表示该slot未被使用

const FSEG_NOT_FULL_N_USED = 8 //FSEG_NOT_FULL链表上被使用的Page数量

const FSEG_FREE = 16 //完全没有被使用并分配给该Segment的Extent链表

const FSEG_NOT_FULL = 16 //至少有一个page分配给当前Segment的Extent链表，全部用完时，转移到FSEG_FULL上，全部释放时，则归还给当前表空间FSP_FREE链表

const FSEG_FULL = 16 //分配给当前segment且Page完全使用完的Extent链表

const FSEG_MAGIC_N = 4 //Magic Number

const FSEG_FRAG_ARR_1 = 4 //属于该Segment的独立Page。总是先从全局分配独立的Page，当填满32个数组项时，就在每次分配时都分配一个完整的Extent，并在XDES PAGE中将其Segment ID设置为当前值

/*--------------------------- Inode Entry struct ---------------------------------------*/
