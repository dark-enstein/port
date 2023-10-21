package config

type Configurer interface {
	// String returns the string representation of all the environment variables
	String() string
	// GetEnv returns the value of a particular env variable
	GetEnv(key string) (interface{}, error)
	// GetEnvs returns an Config containing all the set environmnet variables
	GetEnvs() *Config
}

//type dbConfig struct {
//	Host string `json:"host"`
//	Port string `json:"port"`
//}

type Config struct {
	LogLevel     string      `json:"log_level"`
	Port         string      `json:"server_port"`
	EnabledDB    string      `json:"enabled_db"`
	JWTSecretKey []byte      `json:"jwt_secret_key"`
	DBHost       string      `json:"db_host"`
	Cloud        CloudConfig `json:"cloud"`
}

type CloudConfig struct {
	Provider string `json:"provider"`
	LOC      string `json:"loc"`
}

type EnvBuffer []byte

var env Config

func NewConfig() *Config {
	return &Config{}
}

func NewEnvConfigWithAll(llevel string) *Config {
	return &Config{
		LogLevel: llevel,
	}
}

func (e *Config) ConstructPort() string {
	if e.Port == "" {
		e.Port = "8080"
	}
	return ":" + e.Port
}
