package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

type JSONError struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type JSONSucess struct {
	Msg  interface{} `json:"msg,omitempty"`
	Code int         `json:"code,omitempty"`
	ID   int         `json:"id,omitempty"`
}

func ResponseError(w http.ResponseWriter, status int, message JSONError) {
	RespondJSON(w, status, message)
}

func ResponseSucess(w http.ResponseWriter, sucess JSONSucess) {
	RespondJSON(w, http.StatusOK, sucess)
}

func ResponseBadRequest(w http.ResponseWriter) {
	RespondJSON(w, http.StatusBadRequest, JSONError{Msg: "Bad Request"})
}

func ResponseInternalServerError(w http.ResponseWriter, addMsg ...string) {
	RespondJSON(w, http.StatusInternalServerError, JSONError{Msg: fmt.Sprintf("Internal Server %v", addMsg)})
}
