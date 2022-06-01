package config

type ServerType struct {
	Host                string
	Port                int
	SessionSecretString string
	SessionSecret       []byte `json:"-"`
	Domain              string
	SessionExpire       int
}

type RedisType struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type ZjuOauthType struct {
	ClientID     string
	ClientSecret string
}

type ConfigType struct {
	Server   ServerType
	Redis    RedisType
	ZjuOauth ZjuOauthType
}
