package param

type PagingParam struct {
	PageNo    int // 从1开始
	PageSize  int
	SearchTxt string
}

type SecretParam struct {
	ID       int64  `json:"id"`
	Account  string `json:"account" valid:"required"`
	IsEnable uint8  `json:"isEnable"`
}

type SecretPagingParam struct {
	SecretParam
	PagingParam
}

type ServerParam struct {
	ID                     int64  `json:"id"`
	Sign                   string `json:"sign" valid:"required"`
	Name                   string `json:"name" valid:"required"`
	Remark                 string `json:"remark"`
	IsEnable               uint8  `json:"is_enable"`
	IsOperateSensitiveData uint8  `json:"isOperateSensitiveData"` // 是否可以操作敏感数据（例如：密钥数据），1是，2否
}

type ServerPagingParam struct {
	ServerParam
	PagingParam
}

type GenAccessTokenParam struct {
	ServerSign string `json:"serverSign" valid:"required"`
	TimeToken  string `json:"timeToken" valid:"required"`
}
