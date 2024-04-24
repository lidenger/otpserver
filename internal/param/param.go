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
	Sign                   string `json:"serverSign" valid:"required"`
	Name                   string `json:"serverName" valid:"required"`
	Remark                 string `json:"serverRemark"`
	IsEnable               uint8  `json:"isEnable"`
	IsEnableIPlist         uint8  `json:"isEnableIPlist"`         // 是否启用服务IP白名单，1启用，2禁用
	IsOperateSensitiveData uint8  `json:"isOperateSensitiveData"` // 是否可以操作敏感数据（例如：密钥数据），1是，2否
}

type ServerIpListParam struct {
	Sign string `json:"serverSign" valid:"required"`
	IP   string `json:"ip" valid:"required,ipv4"`
}

type ServerPagingParam struct {
	ServerParam
	PagingParam
}

type GenAccessTokenParam struct {
	ServerSign string `json:"serverSign" valid:"required"`
	TimeToken  string `json:"timeToken" valid:"required"`
}
