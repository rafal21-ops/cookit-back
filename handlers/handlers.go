package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Daniorocket/cookit-back/lib"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Client       *mongo.Client
	DatabaseName string
}
type Pagination struct {
	Data          interface{} `json:"data"`
	Limit         int         `json:"limit"`
	Page          int         `json:"page"`
	TotalElements int         `json:"totalElements"`
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch mux.CurrentRoute(r).GetName() {
		case "CreateRecipe", "ListRecipes", "Renew", "CreateCategory":
			fmt.Println("Auth required!")
			tkn, err := lib.VerifyAndReturnJWT(w, r)
			if err != nil {
				fmt.Println("Error verify token:", err)
				return
			}
			if mux.CurrentRoute(r).GetName() == "Renew" {
				ctx := context.WithValue(r.Context(), "token", tkn)
				rWithCtx := r.WithContext(ctx)
				h.ServeHTTP(w, rWithCtx)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
