package controllers

import (
	"net/http"

	"github.com/arvinpaundra/ecommerce-api/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Hello, my name is Arvin Paundra Ardana!")
}
