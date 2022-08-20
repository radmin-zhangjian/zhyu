package setting

var Server *server
var Database *database
var Redis *redis
var WhiteList *whiteList

type conf struct {
	Srv         server    `yaml:"server"`
	DB          database  `yaml:"database"`
	RedisConfig redis     `yaml:"redis"`
	WhiteList   whiteList `yaml:"whiteList"`
}

type server struct {
	ServerName   string `yaml:"serverName"`
	Port         string `yaml:"port"`
	RunMode      string `yaml:"runMode"`
	LogLevel     string `yaml:"logLevel"`
	LogPath      string `yaml:"logPath"`
	ReadTimeout  int64  `yaml:"readTimeout"`
	WriteTimeout int64  `yaml:"writeTimeout"`
	ShutdownTime int64  `yaml:"shutdownTime"`
	WorkerID     int64  `yaml:"workerID"`
	JwtSecret    string `yaml:"jwtSecret"`
}

type database struct {
	Type            string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	UserName        string `yaml:"username"`
	Password        string `yaml:"password"`
	DbName          string `yaml:"dbname"`
	MaxIdleConn     int64  `yaml:"max_idle_conn"`
	MaxOpenConn     int64  `yaml:"max_open_conn"`
	ConnMaxLifetime int64  `yaml:"conn_max_lifetime"`
}

type redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int64  `yaml:"db"`
	PoolSize int64  `yaml:"poolSize"`
}

type whiteList struct {
	Ip []string `yaml:"ip"`
}
