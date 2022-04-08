package main

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(host, user, password, dbname string) {
    connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

    var err error
    a.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()

    a.initializeRoutes()
}

func (a *App) Run(addr string) {
	fmt.Print("Server is starting at 8080")
    log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) getarticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid article ID")
        return
    }

    p := article{ID: id}
    if err := p.getarticle(a.DB); err != nil {
		fmt.Println("Get Article, err=" + err.Error())
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "article not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }


	resp := getArticleResponse{
		Status: http.StatusOK,
		Message: "Success",
		Data: p,
	}

    respondWithJSON(w, http.StatusOK, resp)
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

func (a *App) getarticles(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    articles, err := getarticles(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

	resp := getArticlesResponse{
		Status: http.StatusOK,
		Message: "Success", 
		Data: articles,
	}
	
    respondWithJSON(w, http.StatusOK, resp)
}

func (a *App) createarticle(w http.ResponseWriter, r *http.Request) {
    var p article
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := p.createarticle(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

	resp := createArticleResponse{
		Status: http.StatusCreated,
		Message: "Success",
		ID: p.ID,
	}
    respondWithJSON(w, http.StatusCreated, resp)
}


func (a *App) updatearticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid article ID")
        return
    }

    var p article
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    p.ID = id

    if err := p.updatearticle(a.DB); err != nil {
		fmt.Println("Update Article, err=" + err.Error())
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deletearticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid article ID")
        return
    }

    p := article{ID: id}
    if err := p.deletearticle(a.DB); err != nil {
		fmt.Println("Delete Article, err=" + err.Error())
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/articles", a.getarticles).Methods("GET")
    a.Router.HandleFunc("/articles", a.createarticle).Methods("POST")
    a.Router.HandleFunc("/articles/{id:[0-9]+}", a.getarticle).Methods("GET")
    a.Router.HandleFunc("/article/{id:[0-9]+}", a.updatearticle).Methods("PUT")
    a.Router.HandleFunc("/article/{id:[0-9]+}", a.deletearticle).Methods("DELETE")
}