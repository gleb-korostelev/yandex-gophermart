package config

import (
	"flag"
	"os"
	"time"

	"github.com/gleb-korostelev/gophermart.git/internal/apperror"
)

const (
	MaxRoutine                  = 20
	DefaultServerAddress        = "localhost:8080"
	DefaultAccuralSystemAddress = "htpp://localhost:8081"
	TokenExpiration             = 24 * time.Hour
	JwtKeySecret                = "very-very-secret-key"
)

type contextKey string

const UserContextKey = contextKey("login")

type ServerConfigData struct {
	ServerAddr           string
	DBDSN                string
	AccuralSystemAddress string
}

var ServerConfig ServerConfigData

func ConfigInit() error {
	flag.StringVar(&ServerConfig.ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&ServerConfig.DBDSN, "d", "", "base file path to save URLs")
	flag.StringVar(&ServerConfig.AccuralSystemAddress, "r", "", "address for accural system")

	flag.Parse()

	if serverAddr := os.Getenv("RUN_ADDRESS"); serverAddr != "" {
		ServerConfig.ServerAddr = serverAddr
	}
	if dbdsn := os.Getenv("DATABASE_URI"); dbdsn != "" {
		ServerConfig.DBDSN = dbdsn
	}
	if accural := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accural != "" {
		ServerConfig.AccuralSystemAddress = accural
	}

	// ServerConfig.DBDSN = "postgres://postgres:7513@localhost:5432/postgres"
	// ServerConfig.ServerAddr = DefaultServerAddress
	// ServerConfig.AccuralSystemAddress = DefaultAccuralSystemAddress

	return checkConfig()
}

func checkConfig() error {
	switch {
	case ServerConfig.ServerAddr == "":
		return apperror.ErrNoServerAddress
	case ServerConfig.DBDSN == "":
		return apperror.ErrNoDatabaseDestination
	case ServerConfig.AccuralSystemAddress == "":
		return apperror.ErrNoAccuralSystemAddress
	default:
		return nil
	}
}
