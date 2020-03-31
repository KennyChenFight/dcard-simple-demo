package model

type Pair struct {
	UserIdOne string `xorm:"pk" json:"userIdOne"`
	UserIdTwo string `xorm:"pk" json:"userIdTwo"`
}

func (p *Pair) TableName() string {
	return "pairs"
}
