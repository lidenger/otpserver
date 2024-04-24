package model

// ServerModel 接入服务
type ServerModel struct {
	ID                     int64  `gorm:"column:id;primary_key" json:"id"`
	Sign                   string `gorm:"column:server_sign" json:"serverSign"`                           // 服务标识
	Name                   string `gorm:"column:server_name" json:"serverName"`                           // 服务名称
	Remark                 string `gorm:"column:server_remark" json:"serverRemark"`                       // 服务描述
	Secret                 string `gorm:"column:server_secret_cipher" json:"serverSecret"`                // 服务密钥
	IV                     string `gorm:"column:server_secret_iv" json:"iv"`                              // 服务密钥IV
	IsOperateSensitiveData uint8  `gorm:"column:is_operate_sensitive_data" json:"isOperateSensitiveData"` // 是否可以操作敏感数据（例如：密钥数据），1是，2否
	IsEnableIPlist         uint8  `gorm:"column:is_enable_iplist" json:"isEnableIPlist"`                  // 是否启用服务IP白名单，1启用，2禁用
	DataCheck              string `gorm:"column:data_check" json:"dataCheck"`                             // 数据校验
	Common
}

func (s *ServerModel) TableName() string {
	return "otp_server"
}
