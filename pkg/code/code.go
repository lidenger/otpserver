package code

type CODE int32

const (
	Success CODE = 200000 // 正常返回无任何异常

	ParamIllegal   CODE = 400001 // 非法参数
	ParamRepeatAdd CODE = 400002 // 重复添加
	DataNotExists  CODE = 400003 // 数据不存在
	Unauthorized   CODE = 401001 // 权限不足
	ReqOverLimit   CODE = 429001 // 请求超限

	ServerErr     CODE = 500000 // 服务异常
	ServerUnready CODE = 500001 // 服务未准备就绪

	StoreErr       CODE = 500010 // 存储异常
	StoreBackupErr CODE = 500011 // 备存储异常

	CryptoErr   CODE = 500020 // 密码操作相关异常
	EncryptErr  CODE = 500021 // 加密异常
	DecryptErr  CODE = 500022 // 解密异常
	DigestedErr CODE = 500023 // 摘要异常
	GenCodeErr  CODE = 500024 // 生成动态令牌异常

	DataErr                   CODE = 500030 // 数据异常
	AccountSecretDataCheckErr CODE = 500031 // 账号密钥数据校验失败

	UnknownErr CODE = 500999 // 未定义的错误
)
