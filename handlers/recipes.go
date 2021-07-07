package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-back/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/validator.v2"
)

func (d *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("token").(JWT)
	recipe := models.Recipe{}
	recipe.UserID = tkn.Username
	recipe.ID = uuid.NewV4().String()
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(recipe); err != nil {
		log.Println("Failed validation:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := models.CreateRecipe(d.Client, d.DatabaseName, &recipe); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (d *Handler) GetListOfRecipes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	recipes, te, err := models.GetAllRecipes(d.Client, d.DatabaseName, page, limit)
	if err != nil {
		log.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pagination{
		Data:          recipes,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	})
}
func (d *Handler) GetListOfRecipesByTags(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]int)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	recipes, te, err := models.GetAllRecipesByTags(d.Client, d.DatabaseName, data["tags"], page, limit)
	if err != nil {
		log.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pagination{
		Data:          recipes,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	})
}
