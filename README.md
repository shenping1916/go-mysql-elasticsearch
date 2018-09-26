原包链接：`https://github.com/siddontang/go-mysql-elasticsearch`

本包是对 `go-mysql-elasticsearch` 的二次开发

主要支持：
+ 对mysql binlog的二次处理，写入nsq
+ 同一条`binlog`事件可按`nsq topic`拆分成多条进行投递 
+ 加入了对`mysqldump`的`binlog`(历史)进行过滤，从而写入`nsq`的`binlog`只是关心的事件
  

## 目录结构
![Alt text](https://github.com/shenping1916/go-mysql-elasticsearch/blob/master/images/1537927471962.png)

## 约定规则：
+ nsq topic："db_name"."db_tablename"."business"
+             数据库名       表名         业务名
  例如：app.user.auth

## sql：
CREATE TABLE `nsq_topic` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `db_name` varchar(30) NOT NULL COMMENT '数据库名',
    `db_table` varchar(30) NOT NULL COMMENT '表名',
    `business` varchar(30) NOT NULL COMMENT '业务名',
    `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否可用：0-不可用; 1-可用',
    `create_time` datetime NOT NULL COMMENT '创建时间',
    `update_time` datetime COMMENT '更新时间',
    `delete_time` datetime COMMENT '删除时间',
    PRIMARY KEY (`id`),
    INDEX `idx_business` (`business`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='NSQ Topic映射表';

## 使用前
### 一、
修改`github.com/siddontang/go-mysql/canal/rows.go`
![Alt text](https://github.com/shenping1916/go-mysql-elasticsearch/blob/master/images/1537928423942.png)

### 二、
修改`github.com/siddontang/go-mysql/canal/dump.go`
![Alt text](https://github.com/shenping1916/go-mysql-elasticsearch/blob/master/images/1537928575018.png)

### 三、
修改`github.com/siddontang/go-mysql/canal/sync.go`
![Alt text](https://github.com/shenping1916/go-mysql-elasticsearch/blob/master/images/1537928880538.png)

## 如何使用?
+ 启动： `./bin/go-mysql-elasticsearch -config=./etc/river.toml`.


