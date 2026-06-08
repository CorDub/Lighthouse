package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateReport (w http.ResponseWriter, r *http.Request) {
	//receive and decode body
	type Body struct {
		ID uuid.UUID `json:"userId" validate:"required"`
	}
	decodedBody := Body{}

	err := decodeRequestBody(r, &decodedBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode the body of the request", err)
		return
	}

	//validate
	errValidation := validate.Struct(decodedBody)
	if errValidation != nil {
		RespondWithError(w, http.StatusInternalServerError, "Invalid user Id provided", errValidation)
		return
	}

	

}