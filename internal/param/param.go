package param

type PagingParam struct {
	PageNo    int // 从1开始
	PageSize  int
	SearchTxt string
}

type SecretParam struct {
	ID       int64  `json:"id"`
	Account  string `json:"account" valid:"required"`
	IsEnable uint8  `json:"is_enable"`
}

type SecretPagingParam struct {
	SecretParam
	PagingParam
}

type ServerParam struct {
	ID       int64  `json:"id"`
	Sign     string `json:"sign" valid:"required"`
	Name     string `json:"name" valid:"required"`
	Remark   string `json:"remark"`
	IsEnable uint8  `json:"is_enable"`
}

type ServerPagingParam struct {
	ServerParam
	PagingParam
}
