package main

import (
	"database/sql"
	"log"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	testHandler := &Handler{}

	return testHandler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

}


