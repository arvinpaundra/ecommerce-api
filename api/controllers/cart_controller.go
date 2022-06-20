package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/arvinpaundra/ecommerce-api/api/auth"
	"github.com/arvinpaundra/ecommerce-api/api/models"
	"github.com/arvinpaundra/ecommerce-api/api/responses"
	"github.com/gorilla/mux"
)

func (server *Server) AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Customer details not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cart := models.Cart{}
	err = json.Unmarshal(body, &cart)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if uint32(customerId) != cid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	cart.CustomerID = uint32(customerId)

	cart.Prepare()
	err = cart.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cartCreated, err := cart.AddToCart(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, cartCreated.CustomerID))
	responses.JSON(w, http.StatusCreated, cartCreated)
}

func (server *Server) GetCustomerCarts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Customer cart not found"))
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	cart := models.Cart{}

	carts, err := cart.FindCustomerCart(server.DB, uint32(customerId))

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, carts)
}

func (server *Server) DeleteCustomerCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Customer details not found"))
		return
	}

	cartId, err := strconv.ParseUint(vars["cartId"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Cart details not found"))
		return
	}

	cid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if cid != uint32(customerId) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	cart := models.Cart{}
	err = server.DB.Debug().Model(models.Cart{}).Where("id = ? AND customer_id = ?", cartId, customerId).Take(&cart).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Cart details not found"))
		return
	}

	_, err = cart.DeleteCart(server.DB, cartId, uint32(customerId))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", cartId))
	responses.JSON(w, http.StatusNoContent, "")
}
