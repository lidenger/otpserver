package model

// ServerIpListModel 接入服务IP白名单
type ServerIpListModel struct {
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	IP         string `gorm:"column:ip" json:"ip"`                  // IP白名单
	ServerSign string `gorm:"column:server_sign" json:"serverSign"` // 服务标识
	Time
}

func (s *ServerIpListModel) TableName() string {
	return "otp_server_ip_whitelist"
}
