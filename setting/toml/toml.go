package toml

var Server *server
var Database *database
var Redis *redis
var Elastic *elastic
var IpWhite *ipWhite

type conf struct {
	Srv         server   `toml:"server"`
	DB          database `toml:"database"`
	RedisConfig redis    `toml:"redis"`
	ES          elastic  `yaml:"elastic"`
	IpWhite     ipWhite  `toml:"ipWhite"`
}

type server struct {
	ServerName   string `toml:"serverName"`
	Port         string `toml:"port"`
	RunMode      string `toml:"runMode"`
	LogLevel     string `toml:"logLevel"`
	LogPath      string `toml:"logPath"`
	ReadTimeout  int64  `toml:"readTimeout"`
	WriteTimeout int64  `toml:"writeTimeout"`
	ShutdownTime int64  `toml:"shutdownTime"`
	WorkerID     int64  `toml:"workerID"`
	JwtSecret    string `toml:"jwtSecret"`
}

type database struct {
	A databaseEn `toml:"a"`
	B databaseEn `toml:"b"`
}

type databaseEn struct {
	Type            string `toml:"type"`
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	UserName        string `toml:"username"`
	Password        string `toml:"password"`
	DbName          string `toml:"dbname"`
	MaxIdleConn     int64  `toml:"max_idle_conn"`
	MaxOpenConn     int64  `toml:"max_open_conn"`
	ConnMaxLifetime int64  `toml:"conn_max_lifetime"`
}

type redis struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Password string `toml:"password"`
	DB       int64  `toml:"db"`
	PoolSize int64  `toml:"poolSize"`
}

type elastic struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ipWhite struct {
	Ip []string `toml:"ip"`
}
