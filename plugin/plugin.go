package plugin

import (
	"fmt"
	"regexp"
	"strconv"
	"dbsync/models"
	"laoyuegou.com/configkit/nsq"
	"go-mysql-elasticsearch/config"
	"go-mysql-elasticsearch/plugin/utility/base"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var driver = "mysql"

type Plugin struct {
	NsqProd   *nsq.NsqProducer
	DbRead    *gorm.DB
}

func NewPlugin(cfg *config.Config) *Plugin {
	plugin := &Plugin{}

	var err error
	port := strconv.Itoa(cfg.SelfMyPort)
	read := []string {
		cfg.SelfMyUser,
		":",
		cfg.SelfMyPasswd,
		"@tcp(",
		cfg.SelfMyHost,
		":",
		port,
		")/",
		cfg.SelfMyDb,
		"?charset=utf8mb4&parseTime=True&loc=Local",
	}
	read_dsn := base.StringSplice(read)
	plugin.DbRead, err = gorm.Open(driver, read_dsn)
	if err != nil {
		panic(err.Error())
	}

	// logger
	plugin.DbRead.LogMode(true)

	// 初始化nsq生产者
	for _, v := range cfg.Nsqs {
		plugin.NsqProd, err = nsq.NewNsqProducer(v.Nsqwrites)
		if err != nil {
			panic(fmt.Errorf("初始化Nsq生产者失败！原因：%v", err.Error()))
		}
	}

	return plugin
}

func (p *Plugin) Truncated(table string) (string, error) {
	reg, err := regexp.Compile(`(\w+)_\d+`)
	if err != nil {
		return "", err
	}

	if reg.MatchString(table) {
		return reg.FindStringSubmatch(table)[1], nil
	} else {
		return table, nil
	}
}

func (p *Plugin) DbQuery(schema, table string) ([]*models.NsqTopic, error) {
	var nsq_topic []*models.NsqTopic
	if err := p.DbRead.Where("db_name = ? AND db_table = ? AND status = ?", schema, table, 1).Find(&nsq_topic).Error; err != nil {
		return nil, err
	}

	return nsq_topic, nil
}

func (p *Plugin) GetNsqTopic(db_name,db_table,business string) string {
	return db_name + "." + db_table + "." + business
}
