// package restapi provides a REST-inspired API via http server
package restapi

import (
	"encoding/json"
	"io/ioutil"
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

func handleCategoryPost(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	vars := mux.Vars(r)
	category := vars["category"]

	b, err := ioutil.ReadAll(r.Body)

	// see stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	var resources map[string]*json.RawMessage
	var resource dr.Dr

	err = json.Unmarshal(b, &resources)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for id, _ := range resources {

		err = json.Unmarshal(*resources[id], &resource)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if resource.Category != category { //avoid cross end-point permission attacks
			http.Error(w, dr.ErrIllegalCategory.Error()+":"+resource.Category, http.StatusInternalServerError)
			return
		}
		if resource.ID != id { //conflicted id
			http.Error(w, dr.ErrUndefinedID.Error()+": did you mean "+resource.ID+" or "+id+"?", http.StatusInternalServerError)
			return
		}
		err = store.Add(resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func handleIDDelete(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	vars := mux.Vars(r)
	category := vars["category"]
	ID := vars["id"]

	_, err := store.Delete(category, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleIDGet(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	vars := mux.Vars(r)
	category := vars["category"]
	ID := vars["id"]

	resource, err := store.Get(category, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func handleIDPost(w http.ResponseWriter, r *http.Request, store dr.Storage) {
	vars := mux.Vars(r)
	category := vars["category"]
	ID := vars["id"]

	b, err := ioutil.ReadAll(r.Body)

	var resource dr.Dr

	err = json.Unmarshal(b, &resource)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resource.Category != category { //avoid cross end-point permission attacks
		http.Error(w, dr.ErrIllegalCategory.Error()+":"+resource.Category, http.StatusInternalServerError)
		return
	}
	if resource.ID != ID { //conflicted id
		http.Error(w, dr.ErrUndefinedID.Error()+": did you mean "+resource.ID+" or "+ID+"?", http.StatusInternalServerError)
		return
	}
	err = store.Add(resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

/*

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
