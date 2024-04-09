package enum

const (
	NacosConfSource   = "nacos"
	LocalConfSource   = "local"
	DefaultConfSource = "default" // 从代码中的[app.toml]加载配置
	UnknownConfSource = "unknown"
)
