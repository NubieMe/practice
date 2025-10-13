package env

import (
	"os"
	"strconv"
	"strings"
)

var (
	GoEnv string

	AllowOrigins []string
	AllowMethods []string
	Port         uint64
	DBHost       string
	DBPort       uint64
	DBUser       string
	DBPassword   string
	DBName       string

	GCSBucket  string
	DirPath    string
	StorageUrl string

	JWTSecretKey string

	AppName string
	Mode    string
)

func GetEnv() {
	GoEnv = os.Getenv("GO_ENV")

	AllowOrigins = strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")
	AllowMethods = strings.Split(os.Getenv("ALLOW_METHODS"), ",")
	Port = parseToUint(os.Getenv("PORT"), 8080)

	DBHost = os.Getenv("DB_HOST")
	DBPort = parseToUint(os.Getenv("DB_PORT"))
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")

	DirPath = emptyDefault(os.Getenv("DIR_PATH"), "uploads")

	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
}

func parseToUint(val string, def ...uint64) uint64 {
	if val == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	parsed, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return parsed
}

func emptyDefault(str string, def string) string {
	if str == "" {
		return def
	}
	return str
}
