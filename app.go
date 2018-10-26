package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App struct
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize method
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	log.Println(connectionString)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Conectado a " + dbname)
	}
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

// Run method
func (a *App) Run(addr string) {
	log.Println("Server listen port 8080")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// InitializeRoutes method
func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/categories", a.getCategories).Methods("GET")
	a.Router.HandleFunc("/category/{id:[0-9]+}", a.getCategory).Methods("GET")
	a.Router.HandleFunc("/categories", a.createCategory).Methods("POST")
	a.Router.HandleFunc("/category/{id:[0-9]+}", a.deleteCategory).Methods("DELETE")
	a.Router.HandleFunc("/category/{id:[0-9]+}", a.updateCategory).Methods("PUT")
}

func (a *App) getCategories(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	categories, err := getCategories(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (a *App) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	p := Category{ID: id}
	if err := p.getCategory(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Category not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createCategory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form")
	}

	var c Category
	c.Description = r.Form.Get("description")
	c.Created_at = time.Now()
	c.Updated_at = time.Now()

	if err := c.createCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	c := Category{ID: id}
	if err := c.deleteCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}

	description := r.Form.Get("description")

	c := Category{ID: id, Description: description, Updated_at: time.Now()}
	if err := c.updateCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}
