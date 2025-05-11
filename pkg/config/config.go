package config

import (
	"os"
	"strconv"
)

type Cfg struct {
	HttpHost       string
	HttpPort       int
	HttpAddr       string
	GrpcHost       string
	GrpcPort       int
	GrpcAddr       string
	AddTimeout     int
	SubTimeout     int
	MltTimeout     int
	DivTimeout     int
	ComputingPower int
}

func getIntEnv(key string, defValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err == nil {
		return value
	}
	return defValue
}

func getHostEnv(key string, defHost string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defHost
}

func NewCfg() *Cfg {
	return &Cfg{
		HttpHost: getHostEnv("ASYNC_CALCULATOR_HTTP_HOST", "localhost"),
		HttpPort: getIntEnv("ASYNC_CALCULATOR_HTTP_PORT", 8080),
		HttpAddr: getHostEnv("ASYNC_CALCULATOR_HTTP_HOST", "localhost") + ":" + strconv.Itoa(getIntEnv("ASYNC_CALCULATOR_PORT", 8080)),

		GrpcHost: getHostEnv("ASYNC_CALCULATOR_GRPC_HOST", "localhost"),
		GrpcPort: getIntEnv("ASYNC_CALCULATOR_GRPC_PORT", 5000),
		GrpcAddr: getHostEnv("ASYNC_CALCULATOR_GRPC_HOST", "localhost") + ":" + strconv.Itoa(getIntEnv("ASYNC_CALCULATOR_GRPC_PORT", 5000)),

		AddTimeout:     getIntEnv("TIME_ADDITION_MS", 5000),
		SubTimeout:     getIntEnv("TIME_SUBTRACTION_MS", 5000),
		MltTimeout:     getIntEnv("TIME_MULTIPLICATIONS_MS", 5000),
		DivTimeout:     getIntEnv("TIME_DIVISIONS_MS", 5000),
		ComputingPower: getIntEnv("COMPUTING_POWER", 10),
	}
}
