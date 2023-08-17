package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"

)

type Employee struct {
	ID        string
	First_Name string
	Last_Name  string
	Department string
	Salary    int
	Age       int
}

type Store interface {
	Connect() error
	Create(employee *Employee) error
	Clear() error
	Get(id string) (*Employee, error)
	Update(id string, employee *Employee) (*Employee, error)
	Delete(id string) error
	List() ([]Employee, error)
}

var dbPool ConnectionPool
var port string

func init() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	privateIp := os.Getenv("DB_PRIVATE_IP")
	port = os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	if dbUser == "" {
		log.Fatal("Please set DB_USER.")
	}
	if dbPassword == "" {
		log.Fatal("Please set DB_PASS")
	}
	if dbName == "" {
		log.Fatal("Please set DB_NAME")
	}
	
	params := Parameters {
		DBUser: dbUser,
		DBPwd: dbPassword,
		DBName: dbName,
		PrivateIP: privateIp,
	}
	
	dbPool = CreateConnectionPool(params)
}

func main() {
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)

	r := mux.NewRouter()
	r.Use(middleware)

	s := handler.Service{SessionStore: newSessionStore()}

	api := r.PathPrefix("/api").Subrouter()
	api.Use(handler.JsonHeader)
	api.Methods(http.MethodOptions).HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	api.Path("/employee").Methods(http.MethodPost).HandlerFunc(s.Create)
	api.Path("/employee").Methods(http.MethodGet).HandlerFunc(s.List)
	api.Path("/employee").Methods(http.MethodDelete).HandlerFunc(s.Clear)
	api.Path("/employee/{id}").Methods(http.MethodGet).HandlerFunc(s.Get)
	api.Path("/employee/{id}").Methods(http.MethodDelete).HandlerFunc(s.Delete)
	api.Path("/employee/{id}").Methods(http.MethodPatch, http.MethodPut).HandlerFunc(s.Update)
	api.Path("/help").Methods(http.MethodGet, http.MethodGet).HandlerFunc(s.Help)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/app/dist")))
	log.Info().Msg("Listening on port " + port)
	log.Fatal().Err(http.ListenAndServe(":" + port, r)).Msg("Can't start service")
}

