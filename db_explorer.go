package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

type Data struct {
	TableName string
	Name string
}

type Response map[string]interface{}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	testHandler := &Handler{DB: db}

	return testHandler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	switch r.Method {
	case "POST":
		h.Create(w, r)
	case "GET":
		h.Read(w, r)
	case "PUT":
		h.Update(w, r)
	case "DELETE":
		h.Delete(w, r)
	default:
		response, _ := json.Marshal(&Response{
		"error": "unknown table",
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("CREATE:", r.URL.Path)

}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	log.Println("READ:", r.URL.Path)

	var db = h.DB

	tableNames, err := getTableNames(db)
	if err != nil {
		response, _ := json.Marshal(&Response{
			"error": "internal server error",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
	}

	for _, tableName := range tableNames {
		log.Println(tableName)
		query := fmt.Sprintf("SELECT * from %s", tableName)
		result, err := h.DB.Query(query)
		if err != nil {
			log.Println("RESULT:", err)
			return
		}

		var data []Data

		for result.Next(){
			d := Data{}
			if err := result.Scan(&d.Name); err != nil {
				log.Println(err)
			}
			log.Println(d.Name)
			data = append(data, d)

			columns, err := result.Columns()
			if err != nil {
				log.Println("COLUMNS:", err)
			}

			for _, v := range columns {
				log.Println(v)
			}
		}

		result.Close()

		log.Println(data)

	}








}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("UPDATE:", r.URL.Path)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE:", r.URL.Path)

}

func getTableNames(db *sql.DB) ([]string, error) {

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Printf("%#v\n", err)
		return nil, err
	}

	var tableNames []string

	for rows.Next() {
		data := Data{}

		if err := rows.Scan(&data.TableName); err != nil {
			log.Printf("%#v\n", err)
			return nil, err
		}

		tableNames = append(tableNames, data.TableName)
	}


	return tableNames, nil
}


