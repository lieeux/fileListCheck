# fileListCheck
读取文件，调用api检查，分类输出到新的文件


有三个文件，文件内容如下：
rfsDataAfidList.txt里是小文件afid列表(小于等于10M且没过期的文件)
rawAfidList.txt里是大文件afid列表(大于10M且没过期的文件)
expiredAfidList.txt里是过期文件afid列表(过期的文件)

afid格式如下：
1e00000000000231ebe7b8a080785e41952427afb43c19602d9a5f998c021a9085aa700fc31259dfc5649e6a7f5b6cd6e742bd03d4c8c4c6fb0461cebc8df000 

读取这三个文件，通过https://pnode.solarfs.io/api/index.html查询这些afid 是否分类正确(小文件是小文件、小文件是不是seed文件、大文件是大文件、过期文件是过期文件)，输出7个文件，内容分别为：
  - 小文件afid正确列表(小于等于10M且没过期的文件)
  - 小文件afid分类异常列表(大于10M或过期的文件)
  - 大文件afid正确列表(大于10M且没过期的文件)
  - 大文件afid分类异常列表(不大于10M或过期的文件)
  - 过期文件afid正确列表(过期的文件)
  - 过期文件afid异常列表(没过期的文件)
  - seed文件列表(调用curl -X GET "http://118.193.47.85:5143/tn/location/seedid/1e00000000f38f9e271fdf870f24d8c07b32efd24bbbb12380831a1b2283cb31928643e6be1be40fb1ff8973720b29dd3ed9cf1d350f4e0545747a56dcf60a24" -H "accept: application/json"这个api，查询小文件是不是seed文件。如果是seed文件，放入seed文件列表。)

http.GET(http://118.193.47.85:5143/tn/location/seedid/1e00000000f38f9e271fdf870f24d8c07b32efd24bbbb12380831a1b2283cb31928643e6be1be40fb1ff8973720b29dd3ed9cf1d350f4e0545747a56dcf60a24) 
结果是json，需要解析json
