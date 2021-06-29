package routing

import (
	"net/http"

	"github.com/Daniorocket/cookit-back/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	NeedAuth    bool
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func initRoutes(handler handlers.Handler) Routes {
	return Routes{
		Route{
			"ListRecipes",
			"GET",
			"/api/v1/recipes",
			false,
			handler.GetListOfRecipes,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/api/v1/recipes",
			true,
			handler.CreateRecipe,
		},
		Route{
			"Login",
			"POST",
			"/api/v1/login",
			false,
			handler.Login,
		},
		Route{
			"Register",
			"POST",
			"/api/v1/register",
			false,
			handler.Register,
		},
		Route{
			"Renew",
			"GET",
			"/api/v1/renew",
			false,
			handler.Renew,
		},
		Route{
			"ListCategories",
			"GET",
			"/api/v1/category",
			false,
			handler.GetListOfCategories,
		},
		Route{
			"CreateCategory",
			"POST",
			"/api/v1/category",
			false,
			handler.CreateCategory,
		},
		Route{
			"GetCategoryByID",
			"GET",
			"/api/v1/category/{id}",
			false,
			handler.GetCategoryByID,
		},
	}
}
