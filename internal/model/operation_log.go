package model

// OperationLogModel 服务操作日志
type OperationLogModel struct {
	ID      int64  `gorm:"column:id;primary_key" json:"id"`
	Item    string `gorm:"column:item" json:"item"`       // 操作项
	Content string `gorm:"column:content" json:"content"` // 操作内容
	IP      string `gorm:"column:ip" json:"ip"`           // 操作者IP
	Time
}

func (s *OperationLogModel) TableName() string {
	return "otp_operation_log"
}
