

- SCAN
  - ```./innodb_view -f ./runoob_tbl.ibd  -t scan -s 16```

- INode Page
  - ```./innodb_view -f ./runoob_tbl.ibd  -t parse -s 16 -d --page_numbers=2```

- Index Page
  - ```./innodb_view -f ./runoob_tbl.ibd  -t parse -s 16 -d --page_numbers=3```


