package model

type SysConfModel struct {
	ID     int64  `gorm:"column:id;primary_key" json:"id"`
	Key    string `gorm:"column:sys_key" json:"key"`
	Val    string `gorm:"column:sys_val" json:"val"`
	Remark string `gorm:"column:remark" json:"remark"`
	Time
}

func (s *SysConfModel) TableName() string {
	return "otp_sys_conf"
}
