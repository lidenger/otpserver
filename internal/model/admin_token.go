package model

type AdminToken struct {
	Account   string `json:"account"`
	Nonce     string `json:"nonce"`
	Time      int64  `json:"time"`
	CheckData string `json:"checkData"`
}
