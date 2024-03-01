package param

type PagingParam struct {
	PageNo    int // 从1开始
	PageSize  int
	SearchTxt string
}

type SecretParam struct {
	ID       int64  `json:"id"`
	Account  string `json:"account"`
	IsEnable uint8  `json:"is_enable"`
}

type SecretPagingParam struct {
	SecretParam
	PagingParam
}

type ServerParam struct {
	ID       int64  `json:"id"`
	Sign     string `json:"sign"`
	Name     string `json:"name"`
	IsEnable uint8  `json:"is_enable"`
}

type ServerPagingParam struct {
	ServerParam
	PagingParam
}
