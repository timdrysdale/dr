package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/timdrysdale/dr"
)

// RESTful API methods
//
// DELETE, GET, POST, UPDATE  /api/resources/<category>/<id>
// DELETE, GET, POST, UPDATE  /api/resources/<category>
// DELETE, GET                /api/resources/
//

func router(store dr.Storage) *mux.Router {

	var router = mux.NewRouter()

	root := "/api"
	//category := `/{category:[a-zA-Z0-9\-\/]+}`
	//id := `/{id:[a-zA-Z0-9\-\/]+}`
	resources := "/resources"

	//router.HandleFunc(root, handleRoot)

	// on root
	// router.HandleFunc(root+resources, handleDeleteRoot).Methods("DELETE")
	router.HandleFunc(root+resources, func(w http.ResponseWriter, r *http.Request) { handleGetRoot(w, r, store) }).Methods("GET")

	// on a specific category
	//router.HandleFunc(root+resources+category, handleDeleteCategory).Methods("DELETE")
	//router.HandleFunc(root+resources+category, handleGetCategory).Methods("GET")
	//router.HandleFunc(root+resources+category, handlePostCategory).Methods("POST", "UPDATE")

	// on a specific id
	//router.HandleFunc(root+resources+category+id, handleDeleteID).Methods("DELETE")
	//router.HandleFunc(root+resources+category+id, handleGetID).Methods("GET")
	//router.HandleFunc(root+resources+category+id, handlePostID).Methods("POST", "UPDATE")

	// other
	//router.HandleFunc(root+"/healthcheck", handleHealthcheck).Methods("GET")

	return router
}
