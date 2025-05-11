package orchestrator

import (
	"net/http"
)

// Функция отправки ошибки 404
func SendNotFoundError404(w http.ResponseWriter) {
	http.Error(w, "404 not found.", http.StatusNotFound)
}

// Функция отправки ошибки 422
func SendForbiddenError403(w http.ResponseWriter) {
	http.Error(w, "{\"error\":\"Authorization required\"}", http.StatusForbidden)
}

// Функция отправки ошибки 422
func SendInvalidDataError422(w http.ResponseWriter) {
	http.Error(w, "{\"error\":\"Expression is not valid\"}", http.StatusUnprocessableEntity)
}

// Функция отправки ошибки 500
func SendInternalError500(w http.ResponseWriter) {
	http.Error(w, "{\"error\":\"Internal server error\"}", http.StatusInternalServerError)
}
