package main

import (
	"github.com/bbconfhq/mycodes/database"
	"github.com/bbconfhq/mycodes/handlers"
	"github.com/bbconfhq/mycodes/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
	"time"

	"flag"
	"log"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":4000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	dsn := os.Getenv("DB_DSN")

	// Check if sqlite file exists
	_, err := os.Stat(dsn)

	// Create new sqlite file
	if err != nil {
		_, _ = os.Create(dsn)
	}

	db := database.Connect(dsn)
	repository.Initialize(db)

	// Delete expired codes hourly; only one process should do it when prefork is enabled
	if !fiber.IsChild() {
		go func() {
			for {
				deleted, err := repository.Code.DeleteExpired()
				if err != nil {
					log.Printf("cleanup error: %v", err)
				} else if deleted > 0 {
					log.Printf("cleanup: deleted %d expired codes", deleted)
				}
				time.Sleep(time.Hour)
			}
		}()
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	handlers.Initialize(app)

	log.Fatal(app.Listen(*port))
}
