package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/thomas3212/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// define application struct to hold application-wide dependencies for the web app
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {

	// Define a new command line flag for network address/port
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// reads in command-line flag and assigns it to addr, needs to be called otherwise addr will always contain default of :4000
	flag.Parse()

	// logger for writing information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate/log.Ltime)

	// logger for error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// initialize new instance of application struct containing dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// needs to make http.Server struct so that server uses our errorLog logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// wraps sql.Open() and returns connection pool for given dsn
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
