package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	PORT := flag.String("port", ":4000", "HTTP listen and serve address")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app := &Application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *PORT,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	log.Printf("Starting server on %v", *PORT)
	err := srv.ListenAndServe()

	errorLog.Fatal(err)
}
