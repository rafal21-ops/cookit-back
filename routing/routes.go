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
			"/v1/recipes",
			false,
			handler.ListRecipes,
		},
		Route{
			"CreateRecipe",
			"POST",
			"/v1/recipes",
			true,
			handler.CreateRecipe,
		},
		Route{
			"Login",
			"POST",
			"/v1/login",
			false,
			handler.Login,
		},
		Route{
			"Register",
			"POST",
			"/v1/register",
			false,
			handler.Register,
		},
		Route{
			"Renew",
			"GET",
			"/v1/renew",
			false,
			handler.Renew,
		},
		Route{
			"ListCategories",
			"GET",
			"/v1/categories",
			false,
			handler.ListCategories,
		},
		Route{
			"CreateCategory",
			"POST",
			"/v1/categories",
			false,
			handler.CreateCategory,
		},
	}
}
