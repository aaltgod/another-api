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
	Name      string
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

		var output []interface{}

		for result.Next() {
			columnNames, err := result.Columns()
			if err != nil {
				log.Println("COLUMNS:", err)
				return
			}

			data := make([]interface{}, len(columnNames))

			columns, err := result.ColumnTypes()
			if err != nil {
				log.Println("ERROR COLUMNS")
				return
			}

			for i, v := range columns {
				columnType := v.DatabaseTypeName()
				switch columnType {
				case "TEXT", "VARCHAR":
					if nullable, _ := v.Nullable(); nullable {
						data[i] = new(sql.NullString)

						break
					}

					data[i] = new(string)
				case "INT":
					if nullable, _ := v.Nullable(); nullable {
						data[i] = new(sql.NullInt32)

						break
					}

					data[i] = new(int)
				}
			}

			if err := result.Scan(data...); err != nil {
				log.Println(err)
				return
			}

			output = append(output, data)
		}

		result.Close()

		for _, v := range output {
			for _, val := range v.([]interface{}) {
				switch val.(type) {
				case *string:
					valType := val.(*string)
					value := *valType

					log.Printf("%#v", value)
				case *int:
					valType := val.(*int)
					value := *valType

					log.Printf("%#v", value)
				case *sql.NullString:
					valType := val.(*sql.NullString)
					value := *valType

					log.Printf("%#v", value)
				case *sql.NullInt32:
					valType := val.(*sql.NullInt32)
					value := *valType

					log.Printf("%#v", value)
				}
			}
		}
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
