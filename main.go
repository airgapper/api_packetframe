package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"

	"github.com/packetframe/api/internal/db"
	"github.com/packetframe/api/internal/routes"
	"github.com/packetframe/api/internal/validation"
)

// Linker flags
var version = "dev"

var (
	document    = flag.Bool("d", false, "Generate documentation instead of starting the API server")
	listenAddr  = flag.String("l", ":8080", "API listen address")
	postgresDSN = flag.String("p", "host=localhost user=api password=api dbname=api port=5432 sslmode=disable", "Postgres DSN")
)

func main() {
	flag.Parse()

	if version == "dev" {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Running in dev mode")
	}

	if !*document {
		log.Println("Connecting to database")
		database, err := db.Connect(*postgresDSN)
		if err != nil {
			log.Fatal(err)
		}
		routes.Database = database

		app := fiber.New(fiber.Config{DisableStartupMessage: true})
		if version == "dev" {
			log.Debugln("Adding wildcard CORS origin")
			app.Use(cors.New())
		}
		routes.Register(app)

		if err := validation.Register(); err != nil {
			log.Fatal(err)
		}

		log.Printf("Starting API on %s", *listenAddr)
		log.Fatal(app.Listen(*listenAddr))
	} else {
		// Document mode
		fmt.Println(routes.Document())
	}
}
