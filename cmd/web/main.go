package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"todo.khoirulakmal.dev/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	todos         *models.TodoModel
	templateCache map[string]*template.Template
	session       *scs.SessionManager
	formDecode    *form.Decoder
}

func main() {

	addr := flag.String("net", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "todouser:Kaskus12@/todo?parseTime=true", "DB Source")
	flag.Parse()

	errorLog := log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Parse template cache
	tmpl, err := parseTemplate()
	if err != nil {
		errorLog.Printf(err.Error())
	} else {
		infoLog.Printf("Parsing template success!")
	}

	// Initialize new session manager
	session := scs.New()
	session.Store = mysqlstore.New(db)
	session.Lifetime = 2 * time.Hour

	// Initialize form decoder
	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		todos:         &models.TodoModel{DB: db},
		templateCache: tmpl,
		session:       session,
		formDecode:    formDecoder,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  errorLog,
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Check if db is connect succesfully
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}
