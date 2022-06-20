package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arvinpaundra/ecommerce-api/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(DbDriver, DBURL)

	if err != nil {
		fmt.Printf("Cannot connect to %s database.", DbDriver)
		log.Fatal("This is error:", err)
	} else {
		fmt.Printf("Connected to the %s database", DbDriver)
	}

	server.DB.Debug().AutoMigrate(&models.Customer{}, &models.Product{}, &models.Category{}, &models.Cart{}, &models.Checkout{}, &models.Payment{})
	server.Router = mux.NewRouter()
	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 5000")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
