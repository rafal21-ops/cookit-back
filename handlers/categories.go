package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/Daniorocket/cookit-back/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (d *Handler) GetListOfCategories(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	categories, te, err := models.GetAllCategories(d.Client, d.DatabaseName, page, limit)
	if err != nil {
		log.Println("Failed to prepare json describe list of users: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pagination{
		Data:          categories,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	})
}
func (d *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cat := models.Category{}
	for {
		part, err := mr.NextPart()
		if err == io.EOF { //End of multipart data
			break
		}
		if err != nil {
			log.Println("Failed:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if part.FormName() == "file" {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, part); err != nil {
				return
			}
			cat.File.EncodedURL = base64.StdEncoding.EncodeToString(buf.Bytes())
			cat.File.Extension = path.Ext(part.FileName())
			cat.ID = uuid.NewV4().String()
			switch ext := cat.File.Extension; ext {
			case ".jpg", ".JPG", ".png", ".PNG":
			default:
				log.Println("Failed:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		//json
		if part.FormName() == "json" {
			jsonDecoder := json.NewDecoder(part)
			err = jsonDecoder.Decode(&cat)
			if err != nil {
				log.Println("Failed:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	if err := models.CreateCategory(d.Client, d.DatabaseName, cat); err != nil {
		log.Println("Failed:", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (d *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	category, err := models.GetCategoryByID(d.Client, d.DatabaseName, id)
	if err != nil {
		log.Println("Failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(category)
}
