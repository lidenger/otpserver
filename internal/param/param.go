package param

type PagingParam struct {
	PageNo    int    `json:"pageNo"` // 从1开始
	PageSize  int    `json:"pageSize"`
	SearchTxt string `json:"searchTxt"`
}

type SecretParam struct {
	ID       int64  `json:"id"`
	Account  string `json:"account" valid:"required"`
	IsEnable uint8  `json:"isEnable"` // 1是，2否
}

type SecretPagingParam struct {
	IsEnable uint8 `json:"isEnable"`
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
