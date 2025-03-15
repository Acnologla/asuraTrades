package main

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/cache"
	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/acnologla/asuraTrades/internal/adapters/http"
	"github.com/acnologla/asuraTrades/internal/adapters/http/controllers"
	"github.com/acnologla/asuraTrades/internal/adapters/http/websocket"
	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/adapters/postgres/repository"
	"github.com/acnologla/asuraTrades/internal/adapters/token"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := config.LoadConfig()
	context := context.Background()
	if err != nil {
		panic(err)
	}

	// initialize adapters

	jwtAdapter := token.NewJwtTokenService(config.JWTConfig)
	postgresConnection := postgres.NewConnection(context, config.PostgresConfig)
	cacheAdapter := cache.NewLocalCache()
	// initialize repositories

	userRepo := repository.NewUserRepository(postgresConnection)
	itemRepo := repository.NewItemRepository(postgresConnection)
	roosterRepo := repository.NewRoosterRepository(postgresConnection)
	userTxProvider := repository.NewTransactionProvider(postgresConnection.(*pgxpool.Pool))

	// initialize services
	userService := service.NewUserService(userRepo, roosterRepo, itemRepo)
	userTokenService := service.NewUserTokenService(jwtAdapter, userService)
	tradeService := service.NewTradeService(cacheAdapter, userService, userTxProvider)

	// initialize controllers
	userTokenController := controllers.NewUserTokenController(config.HTTPConfig.GenerateTokenPassword, userTokenService)
	websocketController := websocket.NewTradeWebsocket(userTokenService, tradeService, config.Production, config.HTTPConfig.ProductionURL)

	// initialize the http server

	http.CreateAndServe(config.HTTPConfig, userTokenController, websocketController)

}
