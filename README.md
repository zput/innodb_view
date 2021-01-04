
#  Viewing about physical file page structure for mysql innodb engine

[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/zput/innodb_view/blob/master/LICENSE)
[![Github Actions](https://github.com/zput/innodb_view/workflows/CI/badge.svg)](https://github.com/zput/innodb_view/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/zput/innodb_view)](https://goreportcard.com/report/github.com/zput/innodb_view)
[![GoDoc](https://godoc.org/github.com/zput/innodb_view?status.svg)](https://godoc.org/github.com/zput/innodb_view)


#### [中文](README-ZH.md) | English

[1. Features](#1-features)

[2. Background](#2-background)

[3. Prerequisites](#3-prerequisites)

[4. Install](#4-quickstart)

[5. Quick Start](#5-quickstart)

[6. Command-line tool](#6-command-line-tool)

[7. Appendix](#7-appendix)

Innodb_view is a Golang implementation for direct access to MySQL InnoDB storage engine files.
The command line allows you to iterate through all the pages that have been used, analyze the type of each page; analyze the Inode page page composition.
Index page page composition, etc. In addition, this project is useful for learning how MySQL innodb physical pages are composed.

## 1. Features

- It is possible to iterate through all the pages that have been used and analyze the type of each page; 
- Analyze the composition of Inode page pages.
- Analyze the composition of Index page pages, etc.


## 2. Background

When learning the underlying knowledge of Innodb, I would like a tool to print out the structure of the file pages on disk visually to facilitate learning and understanding.


## 3. Prerequisites

- Supported MySQL versions: 5.6, 5.7, 8.0.
- Enable innodb_file_per_table, which will create separate *.ibd files for each table.
- InnoDB file page size is set to 16K.

## 4. Install

```bash
go get -u github.com/zput/innodb_view
```

## 5. Quick Start


- SCAN
  - ```./innodb_view -f ./runoob_tbl.ibd  -t scan -s 16```
  
<details>
  <summary>All pages used can be traversed, analyzing the type of each page</summary>
  
```sh
page number:[0]; page type:[File space header]
page number:[1]; page type:[INSERT buffer bitmap page]
page number:[2]; page type:[INODE page(segment object)]
page number:[3]; page type:[INDEX page]
page number:[0]; page type:[Freshly allocated page]
page number:[0]; page type:[Freshly allocated page]
```
</details>

- INode Page
  - ```./innodb_view -f ./runoob_tbl.ibd  -t parse -s 16 -d --page_numbers=2```
  
<details>
  <summary>Analyze Inode page page composition</summary>
  
```sh
+----------------+--------------------------------------------------------------+----------------+
| POSITION       | NAME                                                         | VALUE          |
+----------------+--------------------------------------------------------------+----------------+
| ************** | **************FILE HEADER**************                      | ************** |
| 0 Byte         | checksum                                                     | 3985986369     |
| 4 Byte         | offset                                                       | 2              |
| 8 Byte         | previous_page                                                | 0              |
| 12 Byte        | next_page                                                    | 0              |
| 16 Byte        | lsn_for_last_page_modeification                              | 12621740       |
| 24 Byte        | page_type                                                    | 3              |
| 26 Byte        | flush_lsn                                                    | 0              |
| 34 Byte        | space_id                                                     | 24             |
| ************** | **************index page:list node(first, end)************** | ************** |
| 38 Byte        | first.pageno                                                 | 4294967295     |
| 42 Byte        | first.offset                                                 | 0              |
| 44 Byte        | last.pageno                                                  | 4294967295     |
| 48 Byte        | last.offset                                                  | 0              |
| ************** | **************index page:entry(0-84)**************           | ************** |
| 50 Byte        | [0].fseg_id                                                  | 1              |
| 58 Byte        | [0].fsegnot_fulln_used                                       | 0              |
| 66 Byte        | [0].fseg_free.flst_len                                       | 0              |
| 70 Byte        | [0].fseg_free.first.pageno                                   | 4294967295     |
| 74 Byte        | [0].fseg_free.first.offset                                   | 0              |
| 76 Byte        | [0].fseg_free.last.pageno                                    | 4294967295     |
| 80 Byte        | [0].fseg_free.last.offset                                    | 0              |
| 82 Byte        | [0].fsegnot_full.flst_len                                    | 0              |
| 86 Byte        | [0].fsegnot_full.first.pageno                                | 4294967295     |
| 90 Byte        | [0].fsegnot_full.first.offset                                | 0              |
| 92 Byte        | [0].fsegnot_full.last.pageno                                 | 4294967295     |
| 96 Byte        | [0].fsegnot_full.last.offset                                 | 0              |
| 98 Byte        | [0].fseg_full.flst_len                                       | 0              |
| 102 Byte       | [0].fseg_full.first.pageno                                   | 4294967295     |
| 106 Byte       | [0].fseg_full.first.offset                                   | 0              |
| 108 Byte       | [0].fseg_full.last.pageno                                    | 4294967295     |
| 112 Byte       | [0].fseg_full.last.offset                                    | 0              |
| 114 Byte       | [0].fseg_magicn                                              | 97937874       |
| 118 Byte       | [0].fseg_fragslice[0]                                        | 3              |
|                |                                                              |                |
|                |                                                              |                |
| 122 Byte       | [1].fseg_id                                                  | 2              |
| 130 Byte       | [1].fsegnot_fulln_used                                       | 0              |
| 138 Byte       | [1].fseg_free.flst_len                                       | 0              |
| 142 Byte       | [1].fseg_free.first.pageno                                   | 4294967295     |
| 146 Byte       | [1].fseg_free.first.offset                                   | 0              |
| 148 Byte       | [1].fseg_free.last.pageno                                    | 4294967295     |
| 152 Byte       | [1].fseg_free.last.offset                                    | 0              |
| 154 Byte       | [1].fsegnot_full.flst_len                                    | 0              |
| 158 Byte       | [1].fsegnot_full.first.pageno                                | 4294967295     |
| 162 Byte       | [1].fsegnot_full.first.offset                                | 0              |
| 164 Byte       | [1].fsegnot_full.last.pageno                                 | 4294967295     |
| 168 Byte       | [1].fsegnot_full.last.offset                                 | 0              |
| 170 Byte       | [1].fseg_full.flst_len                                       | 0              |
| 174 Byte       | [1].fseg_full.first.pageno                                   | 4294967295     |
| 178 Byte       | [1].fseg_full.first.offset                                   | 0              |
| 180 Byte       | [1].fseg_full.last.pageno                                    | 4294967295     |
| 184 Byte       | [1].fseg_full.last.offset                                    | 0              |
| 186 Byte       | [1].fseg_magicn                                              | 97937874       |
|                |                                                              |                |
| ************** | **************FILE TRAILER**************                     | ************** |
| N              | oldstyle_checksum                                            | 3985986369     |
| N+4            | lsn                                                          | 12621740       |
| N+8            |                                                              |                |
+----------------+--------------------------------------------------------------+----------------+
```

</details>
  

- Index Page
  - ```./innodb_view -f ./runoob_tbl.ibd  -t parse -s 16 -d --page_numbers=3```
  
<details>
  <summary>Analyze the composition of the Index page</summary>
  
```sh
+----------------+-----------------------------------------------------+----------------+
| POSITION       | NAME                                                | VALUE          |
+----------------+-----------------------------------------------------+----------------+
| ************** | **************FILE HEADER**************             | ************** |
| 0 Byte         | checksum                                            | 2785856177     |
| 4 Byte         | offset                                              | 3              |
| 8 Byte         | previous_page                                       | 4294967295     |
| 12 Byte        | next_page                                           | 4294967295     |
| 16 Byte        | lsn_for_last_page_modeification                     | 12623545       |
| 24 Byte        | page_type                                           | 17855          |
| 26 Byte        | flush_lsn                                           | 0              |
| 34 Byte        | space_id                                            | 24             |
| ************** | **************index page:index header************** | ************** |
| 38 Byte        | ndirslots                                           | 2              |
| 40 Byte        | heap_top                                            | 152            |
| 42 Byte        | nheap                                               | 32771          |
| 44 Byte        | free                                                | 0              |
| 46 Byte        | garbage                                             | 0              |
| 48 Byte        | last_insert                                         | 128            |
| 50 Byte        | direction                                           | 5              |
| 52 Byte        | ndirection                                          | 0              |
| 54 Byte        | nrecs                                               | 1              |
| 56 Byte        | max_trx_id                                          | 0              |
| 64 Byte        | level                                               | 0              |
| 66 Byte        | index_id                                            | 41             |
| ************** | **************index page:FSEG header**************  | ************** |
| 74 Byte        | leafnode.space_id                                   | 24             |
| 78 Byte        | leafnode.node_position.pageno                       | 2              |
| 82 Byte        | leafnode.node_position.offset                       | 242            |
| 84 Byte        | no_leafnode.space_id                                | 24             |
| 88 Byte        | no_leafnode.node_position.pageno                    | 2              |
| 92 Byte        | no_leafnode.node_position.offset                    | 50             |
| ************** | **************index page:All records**************  | ************** |
| 94 Byte        | recordslice[0].info_flags.save_flag1                | 0              |
| 94 Byte 1bits  | recordslice[0].info_flags.save_flag2                | 0              |
| 94 Byte 2bits  | recordslice[0].info_flags.del_flag                  | 0              |
| 94 Byte 3bits  | recordslice[0].info_flags.min_flag                  | 0              |
| 94 Byte 4bits  | recordslice[0].nowned                               | 1              |
| 95 Byte        | recordslice[0].heapno_is_order                      | 0              |
| 96 Byte 5bits  | recordslice[0].record_type                          | 2              |
| 97 Byte        | recordslice[0].next_record_offset_relative          | 29             |
|                |                                                     |                |
| 99 Byte        | recordslice[1].info_flags.save_flag1                | 0              |
| 99 Byte 1bits  | recordslice[1].info_flags.save_flag2                | 0              |
| 99 Byte 2bits  | recordslice[1].info_flags.del_flag                  | 0              |
| 99 Byte 3bits  | recordslice[1].info_flags.min_flag                  | 0              |
| 99 Byte 4bits  | recordslice[1].nowned                               | 0              |
| 100 Byte       | recordslice[1].heapno_is_order                      | 2              |
| 101 Byte 5bits | recordslice[1].record_type                          | 0              |
| 102 Byte       | recordslice[1].next_record_offset_relative          | -16            |
|                |                                                     |                |
| 104 Byte       | recordslice[2].info_flags.save_flag1                | 0              |
| 104 Byte 1bits | recordslice[2].info_flags.save_flag2                | 0              |
| 104 Byte 2bits | recordslice[2].info_flags.del_flag                  | 0              |
| 104 Byte 3bits | recordslice[2].info_flags.min_flag                  | 0              |
| 104 Byte 4bits | recordslice[2].nowned                               | 2              |
| 105 Byte       | recordslice[2].heapno_is_order                      | 1              |
| 106 Byte 5bits | recordslice[2].record_type                          | 3              |
| 107 Byte       | recordslice[2].next_record_offset_relative          | 0              |
|                |                                                     |                |
| 109 Byte       | pagedirectoryslice[0].directoryslot                 | 112            |
|                |                                                     |                |
| 111 Byte       | pagedirectoryslice[1].directoryslot                 | 99             |
|                |                                                     |                |
| ************** | **************FILE TRAILER**************            | ************** |
| N              | oldstyle_checksum                                   | 2785856177     |
| N+4            | lsn                                                 | 12623545       |
| N+8            |                                                     |                |
+----------------+-----------------------------------------------------+----------------+
```

</details>

## 6. Command-line tool

```sh
➜  test git:(develop) ✗ ./innodb_view -h
Usage of ./innodb_view:
  -d, --debug_mode            debug mode (default:false)
  -f, --file_path string      wait parsing file (default "scan")
  -t, --opertor_type string   operator type:(scan/parse) (default "scan")
  -n, --page_numbers ints     page numbers: all page is [-1]; others is [0,1,...] (default [0])
  -s, --page_size int         page size:(16/32 etc) (default 16)
 ```

## 7. Appendix

welcome pr
