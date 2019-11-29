// package restapi provides a REST-inspired API via http server
package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/timdrysdale/dr"
)

const pageNotFound = "page not found"

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, pageNotFound, http.StatusNotFound)
}

func handleResourcesDelete(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	err := store.Reset()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleResourcesGet(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	// list everything we have, in compact form!
	everything, err := store.Categories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(everything)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func handleCategoryGet(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	vars := mux.Vars(r)
	category := vars["category"]

	categoryList, err := store.List(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(categoryList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func handleCategoryDelete(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	vars := mux.Vars(r)
	category := vars["category"]

	categoryList, err := store.List(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for id, _ := range categoryList {
		_, err = store.Delete(category, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

/*

	vars := mux.Vars(r)
	id := vars["id"]

	output, err := json.Marshal(app.Websocket.Rules[id])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}


func
	router.HandleFunc(root, handleRoot)

	// on root
	router.HandleFunc(root+resources, handleDeleteRoot).Methods("DELETE")
	router.HandleFunc(root+resources, handleGetRoot).Methods("GET")

	// on a specific category
	router.HandleFunc(root+resources+category, handleDeleteCategory).Methods("DELETE")
	router.HandleFunc(root+resources+category, handleGetCategory).Methods("GET")
	router.HandleFunc(root+resources+category, handlePostCategory).Methods("POST", "UPDATE")

	// on a specific id
	router.HandleFunc(root+resources+category+id, handleDeleteID).Methods("DELETE")
	router.HandleFunc(root+resources+category+id, handleGetID).Methods("GET")
	router.HandleFunc(root+resources+category+id, handlePostID).Methods("POST", "UPDATE")

	// other
	router.HandleFunc(root+"/healthcheck", handleHealthcheck).Methods("GET")

	return &router
}
*/
