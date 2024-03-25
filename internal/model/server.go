package model

// ServerModel 接入服务
type ServerModel struct {
	ID        int64  `gorm:"column:id;primary_key" json:"id"`
	Sign      string `gorm:"column:server_sign" json:"serverSign"`            // 服务标识
	Name      string `gorm:"column:server_name" json:"serverName"`            // 服务名称
	Remark    string `gorm:"column:server_remark" json:"serverRemark"`        // 服务描述
	Secret    string `gorm:"column:server_secret_cipher" json:"serverSecret"` // 服务描述
	DataCheck string `gorm:"column:data_check" json:"dataCheck"`              // 数据校验 = HMACSHA256(KEY, server_sign + server_secret_cipher + is_enable)
	Common
}

func (s *ServerModel) TableName() string {
	return "otp_server"
}
