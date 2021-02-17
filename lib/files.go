package lib

type File struct {
	Name      string `json:"name" bson:"name"`
	Extension string `json:"extension" bson:"extension"`
	Path      string `json:"path" bson:"path"`
}
