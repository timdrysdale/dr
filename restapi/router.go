package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/timdrysdale/dr"
)

// RESTful API methods from general to specific
//
// ------  GET  ----  ------  /api/healthcheck
// DELETE  GET  ----  ------  /api/resources/
// DELETE  GET  POST  UPDATE  /api/resources/<category>
// DELETE  GET  POST  UPDATE  /api/resources/<category>/<id>

const pathApi = "/api"
const pathResources = pathApi + "/resources"
const pathCategory = pathResources + `/{category:[a-zA-Z0-9\-\/]+}`
const pathID = pathCategory + `/{id:[a-zA-Z0-9\-\/]+}`
const pathHealthcheck = pathApi + "/healthcheck"

func New(store dr.Storage) *mux.Router {

	var router = mux.NewRouter()

	// on root
	router.HandleFunc("/", handleRoot)

	// on all resources
	router.HandleFunc(pathResources,
		func(w http.ResponseWriter, r *http.Request) {
			handleResourcesDelete(w, r, store)
		}).Methods("DELETE")

	router.HandleFunc(pathResources,
		func(w http.ResponseWriter, r *http.Request) {
			handleResourcesGet(w, r, store)
		}).Methods("GET")

	// on a specific category
	router.HandleFunc(pathCategory,
		func(w http.ResponseWriter, r *http.Request) {
			handleCategoryDelete(w, r, store)
		}).Methods("DELETE")

	router.HandleFunc(pathCategory,
		func(w http.ResponseWriter, r *http.Request) {
			handleCategoryGet(w, r, store)
		}).Methods("GET")

	router.HandleFunc(pathCategory,
		func(w http.ResponseWriter, r *http.Request) {
			handleCategoryPost(w, r, store)
		}).Methods("POST", "UPDATE")

	// on a specific ID
	router.HandleFunc(pathID,
		func(w http.ResponseWriter, r *http.Request) {
			handleIDDelete(w, r, store)
		}).Methods("DELETE")

	router.HandleFunc(pathID,
		func(w http.ResponseWriter, r *http.Request) {
			handleIDGet(w, r, store)
		}).Methods("GET")

	router.HandleFunc(pathID,
		func(w http.ResponseWriter, r *http.Request) {
			handleIDPost(w, r, store)
		}).Methods("POST", "UPDATE")

	// on other
	router.HandleFunc(pathHealthcheck,
		func(w http.ResponseWriter, r *http.Request) {
			handleHealthcheck(w, r, store)
		}).Methods("GET")

	return router
}
