package apigateway

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//using http router, register func will do the routing path registration
func (e Endpoints) Register(r *httprouter.Router) {
	r.Handle("POST", "/v1/login", e.HandleLoginPost)
	r.Handler("GET", "/metrics", promttp.Handler())
}

//each method needs a http handler handlers are registered in the register func
func (e Endpoints) HandleLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decodedLoginReq, err := decodeLoginRequest(e.Ctx, r)
	if err != nil {
		respondError(w, 500, err)
		return
	}
	resp, err := e.LoginEndpoint(e.Ctx, decodedLoginReq.(loginRequest))
	if err != nil {
		respondError(w, 500, err)
		return
	}
	respondSuccess(w, resp.(loginResponse))
}

// respondError in some canonical format.
func respondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err,
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

// respondSuccess in some canonical format.
func respondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}