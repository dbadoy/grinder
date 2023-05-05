package api

import "net/http"

func setMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(http.ErrBodyNotAllowed.Error()))
}

func setInternalServerError(w http.ResponseWriter, detail []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(detail)
}
