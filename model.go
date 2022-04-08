package main

import (
	"database/sql"
	"fmt"
)

type createArticleResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	ID int `json:"id"`
}

type getArticlesResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data []article `json:"data"`
}

type getArticleResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data article `json:"data"`
}

type article struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Author string `json:"author"`
}

func (p *article) getarticle(db *sql.DB) error {
    return db.QueryRow("SELECT title, content, author FROM articles WHERE id=$1", p.ID).Scan(&p.Title, &p.Content, &p.Author)
}

func (p *article) updatearticle(db *sql.DB) error {
    _, err :=
        db.Exec("UPDATE articles SET title=$1, content=$2, author=$3 WHERE id=$4", p.Title, p.Content, p.Author, p.ID)

    return err
}

func (p *article) deletearticle(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM articles WHERE id=$1", p.ID)

    return err
}

func (p *article) createarticle(db *sql.DB) error {
    err := db.QueryRow(
        "INSERT INTO articles(title, content, author) VALUES($1, $2, $3) RETURNING id", p.Title, p.Content, p.Author).Scan(&p.ID)

    if err != nil {
		fmt.Println("createarticle, err=" + err.Error())
        return err
    }

    return nil
}

func getarticles(db *sql.DB, start, count int) ([]article, error) {
	rows, err := db.Query(
        "SELECT * FROM articles LIMIT $1 OFFSET $2",
        count, start)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    articles := []article{}

    for rows.Next() {
        var p article
        if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author); err != nil {
            return nil, err
        }
        articles = append(articles, p)
    }

    return articles, nil
}