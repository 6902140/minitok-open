package config

type MysqlConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MinioConfig struct {
	Host            string
	Port            string
	AccessKeyID     string
	SecretAccessKey string
	Videobuckets    string
	Picbuckets      string
}

type RedisConfig struct {
	Host    string
	Port    string
	Network string
	Auth    string
}

type PathConfig struct {
	Videofile string
	Logfile   string
	Picfile   string
}

type Configs struct {
	Mysql MysqlConfig
	Minio MinioConfig
	Redis RedisConfig
	Path  PathConfig
	Level string
}

var Config Configs
