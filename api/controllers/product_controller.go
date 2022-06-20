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
	"github.com/arvinpaundra/ecommerce-api/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) AddProduct(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	product := models.Product{}
	err = json.Unmarshal(body, &product)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	product.Prepare()
	err = product.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productCreated, err := product.SaveProduct(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, productCreated.ID))
	responses.JSON(w, http.StatusCreated, productCreated)
}

func (server *Server) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("category")
	categoryId, err := strconv.ParseUint(query, 10, 32)

	if categoryId != 0 {
		product := models.Product{}
		products, err := product.FindProductByCategory(server.DB, uint32(categoryId))

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		responses.JSON(w, http.StatusOK, products)
		return
	}

	product := models.Product{}
	products, err := product.FindAllProducts(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, products)
}

func (server *Server) GetSingleProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	product := models.Product{}

	productReceived, err := product.FindProductByID(server.DB, productId)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, productReceived)
}

func (server *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	productId, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", productId).Take(&product).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Product not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdate := models.Product{}
	err = json.Unmarshal(body, &productUpdate)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdate.Prepare()
	err = productUpdate.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdate.ID = product.ID

	productUpdated, err := productUpdate.UpdateProduct(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, productUpdated)
}

func (server *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	productId, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", productId).Take(&product).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Product not found"))
		return
	}

	_, err = product.DeleteProduct(server.DB, productId)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", productId))
	responses.JSON(w, http.StatusNoContent, "")
}
