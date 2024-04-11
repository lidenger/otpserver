package config

type M struct {
	Server struct {
		Env                  string `toml:"env"`
		Port                 int    `toml:"port"`
		RootPath             string `toml:"rootPath"`
		ReqLimit             int    `toml:"reqLimit"`
		AccessTokenValidHour int    `toml:"accessTokenValidHour"`
		TimeTokenValidMinute int    `toml:"timeTokenValidMinute"`
		IsEnableLocalStore   bool   `toml:"isEnableLocalStore"`
		IsEnableMemoryStore  bool   `toml:"isEnableMemoryStore"`
	} `toml:"server"`
	Log struct {
		Level      string `toml:"level"`
		RootPath   string `toml:"rootPath"`
		AppFile    string `toml:"appFile"`
		HttpFile   string `toml:"httpFile"`
		MaxSize    int    `toml:"maxSize"`
		MaxBackups int    `toml:"maxBackups"`
		MaxAge     int    `toml:"maxAge"`
		Compress   bool   `toml:"compress"`
	} `toml:"log"`
	Store struct {
		MainStore   string `toml:"mainStore"`
		BackupStore string `toml:"backupStore"`
	} `toml:"store"`
	MySQL struct {
		Address         string `toml:"address"`
		UserName        string `toml:"userName"`
		Password        string `toml:"password@cipher"`
		DbName          string `toml:"dbName"`
		ConnMaxLifeTime int    `toml:"connMaxLifeTime"`
		MaxIdleConn     int    `toml:"maxIdleConn"`
		MaxOpenConn     int    `toml:"maxOpenConn"`
		ConnMaxWaitTime string `toml:"connMaxWaitTime"`
	} `toml:"mysql"`
}

type NacosM struct {
	Client struct {
		NamespaceId string `toml:"namespaceId"`
		DataId      string `toml:"dataId"`
		Group       string `toml:"group"`
		TimeoutMs   uint64 `toml:"timeoutMs"`
		LogDir      string `toml:"logDir"`
		CacheDir    string `toml:"cacheDir"`
		LogLevel    string `toml:"logLevel"`
	} `toml:"client"`
	ServerArr []*NacosServer `toml:"server"`
}

type NacosServer struct {
	Ip          string `toml:"ip"`
	Port        uint64 `toml:"port"`
	ContextPath string `toml:"contextPath"`
}
