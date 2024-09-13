package app

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const GrpcDefaultPort = "6030"

type LKSAPIConfig struct {
	WorkerPoolSize int `validate:"required"`
	DebugMode      bool
}

type FlightInfoAPIConfig struct {
	MaxTabCount int `validate:"required"`
	DebugMode   bool
}

type DBConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
}

type GRPCConfig struct {
	Port string `validate:"required"`
}

type Config struct {
	Database  DBConfig
	GRPC      GRPCConfig
	LKSApi    LKSAPIConfig
	FlightAPI FlightInfoAPIConfig
}

func InitConfig() (Config, error) {
	cfg := Config{}

	path := os.Getenv("CONFIG_FILE_PATH")
	fmt.Println(path)

	if path != "" {
		if err := godotenv.Load(path); err != nil {
			return cfg, fmt.Errorf("load evn's from file %s error: %e", path, err)
		}
	}

	cfg.initDBConfig()
	cfg.initGRPCConfig()
	if err := cfg.initLKSApiConfig(); err != nil {
		return cfg, fmt.Errorf("init lks api config error: %e", err)
	}
	if err := cfg.initFlightInfoApiConfig(); err != nil {
		return cfg, fmt.Errorf("init flight info api config error: %e", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(cfg); err != nil {
		return cfg, fmt.Errorf("validate config error: %e", err)
	}

	return cfg, nil
}

func (c *Config) initDBConfig() {
	c.Database = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func (c *Config) initGRPCConfig() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = GrpcDefaultPort
	} else {
		c.GRPC = GRPCConfig{
			Port: port,
		}
	}
}

func (c *Config) initLKSApiConfig() error {
	pSizeStr := os.Getenv("LKS_API_WORKER_POOL_SIZE")
	pSize, err := strconv.Atoi(pSizeStr)
	if err != nil {
		return fmt.Errorf("LKS_API_WORKER_POOL_SIZE parse error: %e", err)
	}

	debugStr := os.Getenv("LKS_API_DEBUG_MODE")
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		return fmt.Errorf("LKS_API_DEBUG_MODE parse error: %e", err)
	}

	c.LKSApi = LKSAPIConfig{
		WorkerPoolSize: pSize,
		DebugMode:      debug,
	}

	return nil
}

func (c *Config) initFlightInfoApiConfig() error {
	tCountStr := os.Getenv("FLIGHT_INFO_API_TAB_COUNT")
	tCount, err := strconv.Atoi(tCountStr)
	if err != nil {
		return fmt.Errorf("FLIGHT_INFO_API_TAB_COUNT parse error: %e", err)
	}

	debugStr := os.Getenv("FLIGHT_INFO_API_DEBUG_MODE")
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		return fmt.Errorf("FLIGHT_INFO_API_DEBUG_MODE parse error: %e", err)
	}

	c.FlightAPI = FlightInfoAPIConfig{
		MaxTabCount: tCount,
		DebugMode:   debug,
	}

	return nil
}
