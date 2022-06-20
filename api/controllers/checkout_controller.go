package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/arvinpaundra/ecommerce-api/api/auth"
	"github.com/arvinpaundra/ecommerce-api/api/models"
	"github.com/arvinpaundra/ecommerce-api/api/responses"
	"github.com/gorilla/mux"
)

func (server *Server) CustomerCreateCheckout(w http.ResponseWriter, r *http.Request) {
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

	checkout := models.Checkout{}
	err = json.Unmarshal(body, &checkout)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if uint32(customerId) != cid {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	checkout.CustomerID = uint32(customerId)

	checkout.Prepare()
	err = checkout.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	product := models.Product{}
	err = product.BeforeUpdateProudct(server.DB, checkout.Cart.ProductID, checkout.Cart.Qty)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	checkoutCreated, err := checkout.CreateCustomerCheckout(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = product.AfterUpdateProduct(server.DB, checkout.Cart.ProductID, checkout.Cart.Qty)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusCreated, checkoutCreated)
}

func (server *Server) GetCustomerCheckouts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Customer details not found"))
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

	checkout := models.Checkout{}
	checkouts, err := checkout.FindCustomerCheckout(server.DB, customerId)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, checkouts)
}
