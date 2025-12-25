package v1

import (
	"encoding/json"
	"net/http"

	"github.com/solumD/tasks-service/internal/handler/v1/dto"
)

const (
	contentTypeEmpty = ""
	contentTypeJSON  = "application/json"
)

type handler struct {
	taskUsecase TaskUsecase
}

func NewHandler(taskUsecase TaskUsecase) *handler {
	return &handler{
		taskUsecase: taskUsecase,
	}
}

func (h *handler) response(w http.ResponseWriter, contentType string, statusCode int, body []byte) {
	if len(contentType) > 0 {
		w.Header().Add("Content-Type", contentType)
	}

	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	w.Write(body)
}

func (h *handler) errorResponse(w http.ResponseWriter, contentType string, statusCode int, err error) {
	body, errMarsh := json.Marshal(dto.NewErrorResponse(err.Error()))
	if errMarsh != nil {
		h.response(w, contentTypeEmpty, http.StatusInternalServerError, nil)
		return
	}

	if len(contentType) > 0 {
		w.Header().Add("Content-Type", contentType)
	}

	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	w.Write(body)
}
