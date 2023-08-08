package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// define application struct to hold application-wide dependencies for the web app
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Define a new command line flag for network address/port
	addr := flag.String("addr", ":4000", "HTTP network address")

	// reads in command-line flag and assigns it to addr, needs to be called otherwise addr will always contain default of :4000
	flag.Parse()

	// logger for writing information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate/log.Ltime)

	// logger for error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize new instance of application struct containing dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register the file server as the handler for all URL paths that start with "/static/"
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// needs to make http.Server struct so that server uses our errorLog logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
