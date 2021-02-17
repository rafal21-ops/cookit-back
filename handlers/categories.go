package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/Daniorocket/cookit-back/models"
	uuid "github.com/satori/go.uuid"
)

func (d *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories(d.Client, d.DatabaseName)
	if err != nil {
		log.Println("Failed to prepare json describe list of users: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pagination{
		Data:          categories,
		Limit:         1,
		Page:          1,
		TotalElements: len(categories),
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
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Failed:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// file
		if part.FormName() == "file" {
			cat.File.Name = part.FileName()
			cat.File.Path = r.URL.Path + string(os.PathSeparator) + cat.File.Name
			cat.File.Extension = path.Ext(part.FileName())
			switch ext := cat.File.Extension; ext {
			case ".jpg", ".JPG", ".png", ".PNG":
			default:
				log.Println("Failed:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			outfile, err := os.Create("./files/categories/" + part.FileName())
			if err != nil {
				log.Println("Failed:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer outfile.Close()

			_, err = io.Copy(outfile, part)
			if err != nil {
				log.Println("Failed:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
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
	cat.ID = uuid.NewV4().String()
	if err := models.CreateCategory(d.Client, d.DatabaseName, cat); err != nil {
		log.Println("Failed:", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
