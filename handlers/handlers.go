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
type JWT struct {
	Token    string
	Username string
}
type Pagination struct {
	Data          interface{} `json:"data"`
	Limit         string      `json:"limit"`
	Page          string      `json:"page"`
	TotalElements int64       `json:"totalElements"`
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch mux.CurrentRoute(r).GetName() {
		case "CreateRecipe", "ListRecipes", "Renew", "CreateCategory": //Todo
			tkn, err := lib.VerifyAndReturnJWT(w, r)
			if err != nil {
				fmt.Println("Error verify token:", err)
				return
			}
			username, err := lib.GetUsernameFromJWT(tkn)
			if err != nil {
				fmt.Println("Error get username from token:", err)
				return
			}
			ctx := context.WithValue(r.Context(), "token", JWT{Token: tkn, Username: username})
			rWithCtx := r.WithContext(ctx)
			h.ServeHTTP(w, rWithCtx)
			return
		}
		h.ServeHTTP(w, r)
	})
}
