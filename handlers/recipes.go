package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-back/models"
	uuid "github.com/satori/go.uuid"
)

func (d *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("token").(JWT)
	recipe := &models.Recipe{}
	recipe.UserID = tkn.Username
	recipe.ID = uuid.NewV4().String()
	if err := json.NewDecoder(r.Body).Decode(recipe); err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := models.CreateRecipe(d.Client, d.DatabaseName, recipe); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (d *Handler) GetListOfRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := models.GetAllRecipes(d.Client, d.DatabaseName)
	if err != nil {
		log.Println("Failed to prepare json describe list of users: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pagination{
		Data:          recipes,
		Limit:         "1",
		Page:          "1",
		TotalElements: 0,
	})
}
