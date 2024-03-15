package code

type CODE int32

const (
	Success CODE = 200000 // 正常返回无任何异常

	ParamIllegal CODE = 400001 // 参数不合翻

	ServerUnready CODE = 500001 // 服务未就绪
	StoreErr      CODE = 500002 // 存储异常
	UnknownErr    CODE = 500999 // 未定义的错误
)
