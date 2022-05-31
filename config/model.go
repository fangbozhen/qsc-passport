package config

type ServerType struct {
	Host                string
	Port                int
	SessionSecretString string
	SessionSecret       []byte `json:"-"`
}

type MongoDBType struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

type RedisType struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type ConfigType struct {
	Server  ServerType
	MongoDB MongoDBType
	Redis   RedisType
}
