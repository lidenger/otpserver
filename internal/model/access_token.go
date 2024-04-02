package model

type AccessToken struct {
	Sign       string `json:"sign"` // 服务标识
	CreateTime int64  `json:"time"` // 生成时间 time.Now().Unix()
	Rn         string `json:"rn"`   // 随机数
}
