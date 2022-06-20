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

func (server *Server) AddPayment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	payment := models.Payment{}
	err = json.Unmarshal(body, &payment)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	payment.Prepare()
	err = payment.Validate()

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusBadRequest, formattedError)
	}

	paymentCreated, err := payment.SavePayment(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, paymentCreated.ID))
	responses.JSON(w, http.StatusCreated, paymentCreated)
}

func (server *Server) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	payment := models.Payment{}

	payments, err := payment.FindAllPaymentsMethod(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, payments)
}

func (server *Server) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentId, err := strconv.ParseUint(vars["id"], 10, 8)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	payment := models.Payment{}
	err = server.DB.Debug().Model(models.Payment{}).Where("id = ?", paymentId).Take(&payment).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Payment details not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	paymentUpdate := models.Payment{}
	err = json.Unmarshal(body, &paymentUpdate)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	paymentUpdate.Prepare()
	err = paymentUpdate.Validate()

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	paymentUpdate.ID = payment.ID

	paymentUpdated, err := paymentUpdate.UpdatePayment(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, paymentUpdated.ID))
	responses.JSON(w, http.StatusOK, paymentUpdated)
}

func (server *Server) DeletePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentId, err := strconv.ParseUint(vars["id"], 10, 8)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	payment := models.Payment{}
	err = server.DB.Debug().Model(models.Payment{}).Where("id = ?", uint8(paymentId)).Take(&payment).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Payment not found"))
		return
	}

	_, err = payment.DeletePayment(server.DB, uint8(paymentId))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", paymentId))
	responses.JSON(w, http.StatusNoContent, "")
}
