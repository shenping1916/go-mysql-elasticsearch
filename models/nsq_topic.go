package models

import "time"

type NsqTopic struct {
	ID               int64           `gorm:"primary_key" json:"id"`
	DbName           string          `gorm:"column:db_name; type:varchar(30); not null" json:"db_name"`
	DbTable          string          `gorm:"column:db_table; type:varchar(30); not null" json:"db_table"`
	Business         string          `gorm:"column:business; type:varchar(30); index:idx_business; not null" json:"business"`
	Status           int             `gorm:"column:status; type:tinyint(1); not null" json:"status"`
	CreatedAt        time.Time       `gorm:"column:create_time; type:datetime; not null;" json:"create_time"`
	UpdatedAt        time.Time       `gorm:"column:update_time; type:datetime" json:"update_time,omitempty"`
	DeletedAt        *time.Time      `gorm:"column:delete_time; type:datetime" json:"delete_time,omitempty"`
}

func (NsqTopic) TableName() string {
	return "nsq_topic"
}
