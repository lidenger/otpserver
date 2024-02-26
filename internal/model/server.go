package model

// ServerModel 接入服务
type ServerModel struct {
	ID     int64  `gorm:"column:id;primary_key" json:"id"`
	Sign   string `gorm:"column:server_sign" json:"server_sign"`     // 服务标识
	Name   string `gorm:"column:server_name" json:"server_name"`     // 服务名称
	Remark string `gorm:"column:server_remark" json:"server_remark"` // 服务描述
	Common
}

func (s *ServerModel) TableName() string {
	return "otp_server"
}
