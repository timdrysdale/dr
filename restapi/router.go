package restapi

import (
	"fmt"
	"strings"

	"github.com/gorilla/mux"
)

// RESTful API methods
//
// DELETE, GET, POST, UPDATE  $prefix/api/resources/<category>/<id>
// DELETE, GET, POST, UPDATE  $prefix/api/resources/<category>
// DELETE, GET                $prefix/api/resources/
//

func router(prefix string) *mux.Router {

	var router = mux.NewRouter()

	root := slashify(stripslash(prefix) + "/api") // ensure prefix of "/" is properly handled
	category := `/{category:[a-zA-Z0-9\-\/]+}`
	id := `/{id:[a-zA-Z0-9\-\/]+}`
	resources := "/resources"

	router.HandleFunc(root, app.handleRoot)

	// all
	router.HandleFunc(root+resources, handleDeleteAll).Methods("DELETE")
	router.HandleFunc(root+resources, handleGetAll).Methods("GET")

	// on a specific category
	router.HandleFunc(root+resources+category,handleDeleteCategory).Methods("DELETE")
	router.HandleFunc(root+resources+category,handleGetCategory).Methods("GET")
	router.HandleFunc(root+resources+category,handlePostCategory).Methods("POST","UPDATE")

	// on a specific id
	router.HandleFunc(root+resources+category,handleDeleteCategory).Methods("DELETE")
	router.HandleFunc(root+resources+category,handleGetCategory).Methods("GET")
	router.HandleFunc(root+resources+category,handlePostCategory).Methods("POST","UPDATE")

 app.handleDestinationDelete).Methods("DELETE")
	router.HandleFunc("/api/destinations/all", app.handleDestinationShowAll).Methods("GET")
	router.HandleFunc("/api/destinations/all", app.handleDestinationDeleteAll).Methods("DELETE")
	router.HandleFunc(`/api/destinations/{id:[a-zA-Z0-9\-\/]+}`, app.handleDestinationShow).Methods("GET")
	router.HandleFunc("/api/streams", app.handleStreamAdd).Methods("PUT", "POST", "UPDATE")
	router.HandleFunc(`/api/streams/{stream:[a-zA-Z0-9\-\/]+}`, app.handleStreamDelete).Methods("DELETE")
	router.HandleFunc("/api/streams/all", app.handleStreamShowAll).Methods("GET")
	router.HandleFunc("/api/streams/all", app.handleStreamDeleteAll).Methods("DELETE")
	router.HandleFunc(`/api/streams/{stream:[a-zA-Z0-9\-\/]+}`, app.handleStreamShow).Methods("GET")
	router.HandleFunc("/healthcheck", app.handleHealthcheck).Methods("GET")
	router.HandleFunc(`/ts/{feed:[a-zA-Z0-9\-\/]+}`, app.handleTs)

	return &router
}

func stripslash(path string) string {
	path = strings.TrimPrefix(path, "/")
}

func slashify(path string) string {

	//remove trailing slash (that's for directories)
	path = strings.TrimSuffix(path, "/")

	//ensure leading slash
	path = strings.TrimPrefix(path, "/")
	path = fmt.Sprintf("/%s", path)

	return path
}
