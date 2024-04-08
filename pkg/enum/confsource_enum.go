package enum

const (
	ConfSourceNacos   = "nacos"
	ConfSourceLocal   = "local"
	ConfSourceDefault = "default" // 从代码中的[app.toml]加载配置
	ConfSourceUnknown = "unknown"
)
