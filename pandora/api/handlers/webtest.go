package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
)

func WebTest(r chi.Router) {
	r.Mount("/webtest", webtestRouter())
}

func webtestRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getWebTest)
	r.Post("/delete", deleteWebTest)
	r.Put("/", updateWebTest)
	r.Get("/order", saveWebTestOrder)

	return r
}

func getWebTest(w http.ResponseWriter, r *http.Request) {
	// Get the webtest from the database
	var res []models.WebTest
	_ = cache.GetList(constant.PrefixWebTest, &res)

	var order []models.WebTest
	_ = cache.GetList(constant.WebTestOrder, &order)

	// If the order is empty, return the webtest as is
	if len(order) == 0 {
		render.JSON(w, r, res)
		return
	}

	// Create a map for quick lookup of webtest by ID
	webtestMap := make(map[string]models.WebTest)
	for _, webtest := range res {
		webtestMap[webtest.Id] = webtest
	}

	// Sort res based on the order
	var sortedRes []models.WebTest
	for _, item := range order {
		if webtest, exists := webtestMap[item.Id]; exists {
			sortedRes = append(sortedRes, webtest)
		}
	}

	render.JSON(w, r, sortedRes)
}

func deleteWebTest(w http.ResponseWriter, r *http.Request) {
	webtest := &models.WebTest{}
	if err := render.DecodeJSON(r.Body, webtest); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// Delete the webtest from the database
	_ = cache.Delete(webtest.Id)

	render.NoContent(w, r)
}

func updateWebTest(w http.ResponseWriter, r *http.Request) {
	webtest := &models.WebTest{}
	if err := render.DecodeJSON(r.Body, webtest); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	if webtest.Id == "" {
		webtest.Id = fmt.Sprintf("%s%d", constant.PrefixWebTest, utils.SnowflakeId())
	}

	// Add the webtest to the database
	_ = cache.Put(webtest.Id, webtest)

	render.NoContent(w, r)
}
