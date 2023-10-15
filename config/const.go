package config

var (
	ENABLED int
)

const (
	FLAG   = iota // --flags
	SHELL         // TODO ref: https://docs.docker.com/compose/environment-variables/set-environment-variables/#substitute-from-the-shell
	ENV           // ENV_VARS
	DOTENV        // .env file at root
)
