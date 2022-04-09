package main

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/joho/godotenv"
)

var a App

func TestMain(m *testing.M) {
	godotenv.Load()
    a.Initialize(
		os.Getenv("APP_DB_HOST"),
        os.Getenv("APP_DB_USERNAME"),
        os.Getenv("APP_DB_PASSWORD"),
        os.Getenv("APP_DB_NAME"))

    ensureTableExists()
    code := m.Run()
    clearTable()
    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("DELETE FROM articles")
    a.DB.Exec("ALTER SEQUENCE articles_id_seq RESTART WITH 1")
}

// const tableDropQuery = `DROP TABLE IF EXISTS articles`

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS articles
(
    id SERIAL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
	author TEXT NOT NULL,
    CONSTRAINT articles_pkey PRIMARY KEY (id)
)`


func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/articles", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "{\"status\":200,\"message\":\"Success\",\"data\":[]}" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestGetNonExistentarticle(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/articles/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "article not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'article not found'. Got '%s'", m["error"])
    }
}


func TestCreateArticle(t *testing.T) {

    clearTable()

    var jsonStr = []byte(`{"title":"test article", "content": "content", "author": "author"}`)
    req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response := executeRequest(req)
    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    //fmt.Println(m);

    if m["message"] != "Success" {
        t.Errorf("Expected message to be 'Success'. Got '%v'", m["message"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected article ID to be '1'. Got '%v'", m["id"])
    }
}


func TestGetarticle(t *testing.T) {
    clearTable()
    addarticles(1)

    req, _ := http.NewRequest("GET", "/articles/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addarticles(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
        a.DB.Exec("INSERT INTO articles(title, content, author) VALUES($1, $2, $3)", "title "+strconv.Itoa(i), "content "+strconv.Itoa(i), "author "+strconv.Itoa(i))
    }
}


/*func TestUpdatearticle(t *testing.T) {

    clearTable()
    addarticles(1)

    req, _ := http.NewRequest("GET", "/articles/1", nil)
    response := executeRequest(req)
    var originalarticle map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalarticle)

    var jsonStr = []byte(`{"title":"test article - updated name", "content": "content", "author": "author"}`)
    req, _ = http.NewRequest("PUT", "/article/1", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    fmt.Println(originalarticle)
    fmt.Println(m)

    if m["id"] != originalarticle["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalarticle["id"], m["id"])
    }

    if m["title"] == originalarticle["title"] {
        t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalarticle["title"], m["title"], m["title"])
    }

    if m["content"] == originalarticle["content"] {
        t.Errorf("Expected the content to change from '%v' to '%v'. Got '%v'", originalarticle["content"], m["content"], m["content"])
    }

	if m["author"] == originalarticle["author"] {
        t.Errorf("Expected the author to change from '%v' to '%v'. Got '%v'", originalarticle["author"], m["author"], m["author"])
    }
}*/


/*func TestDeletearticle(t *testing.T) {
    clearTable()
    addarticles(1)

    req, _ := http.NewRequest("GET", "/article/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/article/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/article/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}*/