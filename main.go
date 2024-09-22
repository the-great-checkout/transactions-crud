package main

import (
	"github.com/Netflix/go-env"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/the-great-checkout/transactions-crud/docs"

	"github.com/the-great-checkout/transactions-crud/internal/controller"
	"github.com/the-great-checkout/transactions-crud/internal/database"
	"github.com/the-great-checkout/transactions-crud/internal/mapper"
	"github.com/the-great-checkout/transactions-crud/internal/repository"
	"github.com/the-great-checkout/transactions-crud/internal/service"
)

type Environment struct {
	Postgres struct {
		Schema string `env:"POSTGRES_SCHEMA,default=transactions."`
		DSN    string `env:"POSTGRES_DSN,default=host=localhost user=user password=password dbname=postgres port=5432 sslmode=disable"`
	}
	Mongo struct {
		URI        string `env:"MONGO_URI,default=mongodb://root:example@localhost:27017"`
		Database   string `env:"MONGO_DATABASE,default=transactions-crud"`
		Collection string `env:"MONGO_COLLECTION,default=transactions"`
	}

	Kafka struct {
		Topic   string `env:"KAFKA_TOPIC,default=transactions"`
		Address string `env:"KAFKA_ADDRESS,default=localhost:9092"`
	}

	Port string `env:"PORT,default=:8081"`
}

//	@title			Transactions CRUD API
//	@version		1.0
//	@description	This is a sample server for transactions CRUD.
//	@host			localhost:8081
//	@BasePath		/

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		panic(err)
	}

	mongo := database.NewMongo(environment.Mongo.URI, environment.Mongo.Database, environment.Mongo.Collection)
	postgres := database.NewPostgres(environment.Postgres.DSN, environment.Postgres.Schema)

	notificationService := service.NewNotificationService(environment.Kafka.Topic, environment.Kafka.Address)

	transactionRepository := repository.NewTransactionRepository(postgres, mongo)
	transactionMapper := mapper.NewTransactionMapper()
	transactionService := service.NewTransactionService(transactionRepository, transactionMapper)
	transactionController := controller.NewTransactionController(transactionService, notificationService)

	statusRepository := repository.NewStatusRepository(postgres)
	statusMapper := mapper.NewStatusMapper()
	statusService := service.NewStatusService(statusRepository, statusMapper)
	statusController := controller.NewStatusController(statusService)

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/v1")

	v1.POST("/transactions", transactionController.CreateHandler)
	v1.GET("/transactions/:transactionID", transactionController.GetByIDHandler)
	v1.GET("/transactions", transactionController.GetAllHandler)
	v1.PUT("/transactions/:transactionID", transactionController.UpdateHandler)
	v1.DELETE("/transactions/:transactionID", transactionController.DeleteHandler)
	v1.POST("/statuses", statusController.CreateHandler)
	v1.GET("/statuses/:statusID", statusController.GetByIDHandler)
	v1.GET("/statuses", statusController.GetAllHandler)

	e.Logger.Fatal(e.Start(environment.Port))
}
