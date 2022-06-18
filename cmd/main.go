package main

import (
	"context"
	"log"

	"github.com/romik1505/ApiGateway/internal/app/config"
	"github.com/romik1505/ApiGateway/internal/app/server"
	"github.com/romik1505/ApiGateway/internal/app/service"
)

// @title           Auth service
// @version         0.1
// @description     This is authenticate service.
// @host 			localhost:8080
// @BasePath 		/
// @in header
// @name Authorization

func main() {
	ctx := context.Background()

	mongoConnection := config.NewMongoConnection(ctx, config.GetValue(config.MongoConnection))

	userService := service.NewUserService(ctx, mongoConnection)
	app := server.NewApp(ctx, userService)
	if err := app.Run(); err != nil {
		log.Println(err.Error())
	}
}
