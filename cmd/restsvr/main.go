package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/echohandler"
	"github.com/JHTohru/eulabschallenge/pkg/postgres"
	"github.com/JHTohru/eulabschallenge/pkg/product"
)

var port = flag.String("port", "8080", "rest server port")

func main() {
	// connect to database
	dsn := "host=localhost user=postgres password=password " +
		"dbname=eulabschallengedb port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// initialize data access layer objects
	counter := postgres.NewProductCounter(db)
	fetcher := postgres.NewProductFetcher(db)
	finder := postgres.NewProductFinder(db)
	inserter := postgres.NewProductInserter(db)
	remover := postgres.NewProductRemover(db)
	saver := postgres.NewProductSaver(db)

	// initialize business layer objects
	creator := product.NewCreator(inserter)
	deleter := product.NewDeleter(finder, remover)
	getter := product.NewGetter(finder)
	lister := product.NewLister(counter, fetcher)
	updater := product.NewUpdater(finder, saver)

	// initialize transport layer objects
	createHandler := echohandler.NewCreateHandler(creator)
	deleteHandler := echohandler.NewDeleteHandler(deleter)
	getHandler := echohandler.NewGetHandler(getter)
	listHandler := echohandler.NewListHandler(lister)
	updateHandler := echohandler.NewUpdateHandler(updater)

	// register routes
	e := echo.New()
	e.POST("/products", createHandler.Handle)
	e.GET("/products", listHandler.Handle)
	e.GET("/products/:id", getHandler.Handle)
	e.PUT("/products/:id", updateHandler.Handle)
	e.DELETE("/products/:id", deleteHandler.Handle)
	e.HTTPErrorHandler = echohandler.HTTPErrorHandler

	// start server
	e.Logger.Fatal(e.Start(":" + *port))
}
