package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type jwtOptions struct {
	Secret   string
	Issuer   string
	Audience string
}

type Config struct {
	DSN        string
	JwtOptions jwtOptions
}

var Cfg Config

func Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se pudo cargar .env, usando variables de entorno existentes")
	}

	Cfg = Config{
		DSN: os.Getenv("DB_DSN"),
		JwtOptions: jwtOptions{
			Secret:   os.Getenv("JWT_SECRET"),
			Issuer:   os.Getenv("JWT_ISSUER"),
			Audience: os.Getenv("JWT_AUDIENCE"),
		},
	}
}
