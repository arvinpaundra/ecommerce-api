package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/arvinpaundra/ecommerce-api/api/models"
	"github.com/arvinpaundra/ecommerce-api/api/responses"
	"github.com/arvinpaundra/ecommerce-api/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) AddCategory(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	category := models.Category{}
	err = json.Unmarshal(body, &category)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	category.Prepare()
	err = category.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryCreated, err := category.SaveCategory(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	w.Header().Set("Lcation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, categoryCreated.ID))
	responses.JSON(w, http.StatusCreated, categoryCreated)
}

func (server *Server) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}

	categories, err := category.FindAllCategories(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, categories)
}

func (server *Server) GetSingleCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	category := models.Category{}

	categoryReceived, err := category.FindCategoryByID(server.DB, uint32(categoryId))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, categoryReceived)
}

func (server *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	category := models.Category{}
	err = server.DB.Debug().Model(models.Category{}).Where("id = ?", uint32(categoryId)).Take(&category).Error

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Category not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryUpdate := models.Category{}
	err = json.Unmarshal(body, &categoryUpdate)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryUpdate.Prepare()
	err = categoryUpdate.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryUpdate.ID = category.ID

	categoryUpdated, err := categoryUpdate.UpdateCategory(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, categoryUpdated.ID))
	responses.JSON(w, http.StatusOK, categoryUpdated)
}

func (server *Server) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	category := models.Category{}
	err = server.DB.Debug().Model(&models.Category{}).Where("id = ?", uint32(categoryId)).Take(&category).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Category not found"))
		return
	}

	_, err = category.DeleteCategory(server.DB, uint32(categoryId))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", categoryId))
	responses.JSON(w, http.StatusNoContent, "")
}
