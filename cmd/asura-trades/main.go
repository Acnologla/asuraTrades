package main

import (
	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/acnologla/asuraTrades/internal/adapters/http"
	"github.com/acnologla/asuraTrades/internal/adapters/http/controllers"
	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/adapters/postgres/repository"
	"github.com/acnologla/asuraTrades/internal/adapters/token"
	"github.com/acnologla/asuraTrades/internal/core/service"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// initialize adapters

	jwtAdapter := token.NewJwtTokenService(config.JWTConfig)
	postgresConnection := postgres.New(config.PostgresConfig)

	// initialize repositories

	userRepo := repository.NewUserRepository(postgresConnection)
	itemRepo := repository.NewItemRepository(postgresConnection)
	roosterRepo := repository.NewRoosterRepository(postgresConnection)

	// initialize services
	userService := service.NewUserService(userRepo, roosterRepo, itemRepo)
	userTokenService := service.NewUserTokenService(jwtAdapter, userService)

	// initialize controllers
	userTokenController := controllers.NewUserTokenController(config.HTTPConfig.GenerateTokenPassword, userTokenService)

	// initialize the http server

	http.CreateAndServe(config.HTTPConfig, userTokenController)

}
