package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	DB *sql.DB
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

	log.Println(h.DB.Stats())

	_, err := h.DB.Query("SELECT *")
	if err != nil {
		log.Println("ROWS:", err)
		response, _ := json.Marshal(&Response{
			"error": "internal server error",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)

		return
	}

	//result, err := h.DB.Exec("SELECT *")
	//if err != nil {
	//	log.Println("RESULT:", err)
	//	return
	//}
	//affected, err := result.RowsAffected()
	//if err != nil {
	//	log.Println("AFFECTED", err)
	//	return
	//}
	//log.Println(affected)


}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("UPDATE:", r.URL.Path)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE:", r.URL.Path)

}


