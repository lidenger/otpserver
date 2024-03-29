package model

type AccessToken struct {
	Sign string `json:"sign"` // 服务标识
	Tim  int64  `json:"time"` // 生成时间
	Rn   string `json:"rn"`   // 随机数
}
