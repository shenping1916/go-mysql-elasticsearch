package plugin

import (
	"go-mysql-elasticsearch/elastic"
	"go-mysql-elasticsearch/plugin/utility/base"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
)

func (p *Plugin) Handler(reqs []*elastic.BulkRequest) {
	// 写入nsq
	for _, req := range reqs {
		is_hitory := req.Event.History
		// 过滤mysqldump 的binlog事件
		if !is_hitory {
			// 分离数据库名、表名
			schema := req.Event.Table.Schema
			table, err := p.Truncated(req.Event.Table.Name)
			if err != nil {
				log.Errorf("正则匹配出错！%v", err.Error())
				return
			}

			// 获取nsq_topic查询结果集
			records, err := p.DbQuery(schema, table)
			if err != nil {
				log.Error("数据库查询错误！%w", err.Error())
				return
			}

			if len(records) == 0 {
				//log.Warn("Table: [nsq_topic] Record Not Found !")
				return
			}

			for _, record := range records {
				topic := p.GetNsqTopic(record.DbName, record.DbTable, record.Business)

				switch req.Event.Action {
				case canal.UpdateAction:
					smap := &base.Smap{}
					data := make([]map[string]interface{}, 0, len(req.Event.Rows))
					for _, row := range req.Event.Rows {
						r := row
						m := make(map[string]interface{})
						for i := 0; i < len(req.Event.Table.Columns); i++ {
							field := req.Event.Table.Columns[i].Name
							value := r[i]
							switch value := value.(type) {
							case []uint8:
								m[field] = string(value)
							default:
								m[field] = value
							}
						}

						data = append(data, m)
					}

					// 写入nsq
					smap.Put("topic", topic)
					smap.Put("action", req.Event.Action)
					smap.Put("data_before", data[0])
					smap.Put("data", data[1])
					s := base.ToSortedMapJson(smap)
					s_byte := base.StringToBytes(s)

					if err := p.NsqProd.SendNsqMsg(topic, s_byte); err != nil {
						log.Errorf("binlog消息写入Nsq失败！原因：%v", err.Error())
						return
					}
					log.Infof("binlog消息写入Nsq成功: %s", s)

				case canal.InsertAction, canal.DeleteAction:
					smap := &base.Smap{}
					s1 := &base.Smap{}
					for _, row := range req.Event.Rows {
						r := row
						for i := 0; i < len(req.Event.Table.Columns); i++ {
							field := req.Event.Table.Columns[i].Name
							value := r[i]
							switch value := value.(type) {
							case []uint8:
								s1.Put(field, string(value))
							default:
								s1.Put(field, value)
							}
						}

						// 写入nsq
						smap.Put("topic", topic)
						smap.Put("action", req.Event.Action)
						smap.Put("data", s1)
						s := base.ToSortedMapJson(smap)
						s_byte := base.StringToBytes(s)

						if err := p.NsqProd.SendNsqMsg(topic, s_byte); err != nil {
							log.Errorf("binlog消息写入Nsq失败！原因：%v", err.Error())
							return
						}
						log.Infof("binlog消息写入Nsq成功: %s", s)
					}
				}
			}
		}
	}
}

