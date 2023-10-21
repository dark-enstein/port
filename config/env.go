package config

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
)

const (
	DefaultSecretKey = "jfbireowlqjewdoiuxasukldkewjhfdqjsuewfdikewjkhjoieoiqjwniewrnf"
)

var (
	configJsonTemplate = `[     {"log_level": "%s"}, {"server_port": "%s"},     {"enabled_db": "%s"},    {"db_host": "%s"}, {"cloud": {{"provider": "%s"}, {"loc": "%s"}, {"content": "%s"}} ]`
)

var (
	ErrConfigNotFound = "requested config not found"
)

type fileLoc string

func (f *fileLoc) toString() string {
	ref := reflect.ValueOf(*f)
	if ref.Kind() != reflect.String {
		return ""
	}
	return ref.String()
}

func (f *fileLoc) Read() []byte {
	file, err := os.ReadFile(f.toString())
	if err != nil {
		_ = fmt.Errorf("cannot read File in loc due to err: %w", err)
		return nil
	}
	return file
}

type EnvConfig struct {
	logLevel     string
	port         string
	enabledDB    string
	jwtSecretLoc string
	dBHost       string
	cloud        struct {
		provider string
		loc      string
		content  []byte // AWS: AKey, ASecret, Profile, Region // GCP, etc
	}
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

// TODO: String returns the string representation of the config
func (e *EnvConfig) String() string {
	cfg := fmt.Sprintf(configJsonTemplate, e.logLevel, e.port, e.enabledDB, e.dBHost, e.cloud.provider, e.cloud.loc, e.cloud.content)
	//cfgByte, _ := json.Marshal(cfg) // TODO: think ways to print the config string to any io.Writer. Why?
	//json.Unmarshal()
	return cfg
}

func (e *EnvConfig) flatten() *map[string]interface{} {
	mapConfig := map[string]interface{}{}
	err := mapstructure.Decode(e, &mapConfig)
	if err != nil {
		return nil
	}
	return &mapConfig
}

func (e *EnvConfig) GetEnvs() *Config {
	loc := fileLoc(e.jwtSecretLoc)
	return &Config{
		LogLevel:     e.logLevel,
		Port:         e.port,
		EnabledDB:    e.enabledDB,
		JWTSecretKey: resolveJWTSecretKey(loc),
		DBHost:       e.enabledDB,
		Cloud: CloudConfig{
			Provider: e.cloud.provider,
			LOC:      e.cloud.loc,
		},
	}
}

func resolveJWTSecretKey(loc fileLoc) []byte {
	if loc.Read() == nil {
		return []byte(DefaultSecretKey)
	}
	return loc.Read()
}

func (e *EnvConfig) GetEnv(v string) (interface{}, error) {
	m := *e.flatten()
	val, ok := m[v]
	if !ok {
		return "", errors.New(ErrConfigNotFound)
	}
	return val, nil
}
