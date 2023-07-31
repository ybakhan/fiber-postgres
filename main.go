package main

import (
	"fmt"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-kit/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ybakhan/fiber-postgres/config"
	"github.com/ybakhan/fiber-postgres/controller"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := config.ReadConfig()
	logger := initializeLogger()
	logger.Log("msg", "fiber postgres started", "configuration", &config)

	app := fiber.New()
	api := app.Group(fmt.Sprintf("/%s/books", config.Version))

	DB, err := initializeDb(config)
	if err != nil {
		fmt.Println("failed to create db")
		return
	}

	if err := DB.AutoMigrate(&controller.Book{}); err != nil {
		fmt.Println("failed to migrate db")
		return
	}

	cache := memcache.New(fmt.Sprintf("%s:%d", config.Cache.Host, config.Cache.Port))

	bookController := controller.NewBookController(DB, cache)
	api.Post("", bookController.CreateBook)
	api.Get("/:id", bookController.GetBook)
	api.Get("", bookController.GetBooks)
	api.Delete("/:id", bookController.DeleteBook)

	app.Listen(fmt.Sprintf(":%d", config.Port))
}

func initializeLogger() log.Logger {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

func initializeDb(config *config.Config) (*gorm.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSLMode,
	)

	return gorm.Open(postgres.Open(connStr), &gorm.Config{})
}
