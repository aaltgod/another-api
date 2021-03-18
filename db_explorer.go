package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	DB     *sql.DB
	Output interface{}
}

type DB struct {
	TableName string
	Name      string
}

type Response map[string]interface{}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	testHandler := &Handler{DB: db}

	return testHandler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("URL PATH", r.URL.Path)

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
		internalServerError(w)

		return
	}

	reqTableName := r.URL.Path

	switch reqTableName {
	case "/":
		response, err := json.Marshal(&Response{
			"response": Response{
				"tables": tableNames,
			},
		})
		if err != nil {
			internalServerError(w)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)

	default:
		for _, tableName := range tableNames {
			if tableName == strings.Trim(reqTableName, "/") {
				query := fmt.Sprintf("SELECT * FROM %s ", tableName)

				if isExistsParam(r, "offset") {
					param, err := strconv.Atoi(r.FormValue("offset"))
					if err != nil {
						offset := 0
						query += fmt.Sprintf("WHERE id > %d ", offset)

						break
					}

					offset := param
					query += fmt.Sprintf("WHERE id > %d ", offset)
				}

				if isExistsParam(r, "limit") {
					param, err := strconv.Atoi(r.FormValue("limit"))
					if err != nil {
						limit := 5
						query += fmt.Sprintf("LIMIT %d", limit)

						break
					}

					limit := param
					query += fmt.Sprintf("LIMIT %d", limit)
				}

				result, err := h.DB.Query(query)
				if err != nil {
					log.Println("RESULT:", err)
					return
				}

				var output []interface{}

				for result.Next() {
					dataMap := make(map[string]interface{})
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
								dataMap[v.Name()] = data[i]

								break
							}

							data[i] = new(string)
							dataMap[v.Name()] = data[i]
						case "INT":
							if nullable, _ := v.Nullable(); nullable {
								data[i] = new(sql.NullInt32)
								dataMap[v.Name()] = data[i]

								break
							}

							data[i] = new(int)
							dataMap[v.Name()] = data[i]
						}
					}

					if err := result.Scan(data...); err != nil {
						log.Println(err)
						return
					}

					output = append(output, dataMap)
				}

				result.Close()

				for _, v := range output {
					oneSet := v.(map[string]interface{})
					for key, val := range oneSet {
						switch val.(type) {
						case *string:
							valType := val.(*string)
							value := *valType
							oneSet[key] = value
						case *int:
							valType := val.(*int)
							value := *valType
							oneSet[key] = value
						case *sql.NullString:
							valType := val.(*sql.NullString)
							value := *valType
							if value.Valid {
								oneSet[key] = value.String

								break
							}

							oneSet[key] = nil
						case *sql.NullInt32:
							valType := val.(*sql.NullInt32)
							value := *valType
							if value.Valid {
								oneSet[key] = value.Int32

								break
							}

							oneSet[key] = nil
						}
					}
				}

				response, err := json.Marshal(&Response{
					"response": Response{
						"records": output,
					},
				})
				if err != nil {
					log.Println(err)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(response)

				return
			}
		}

		response, _ := json.Marshal(&Response{
			"error": "unknown table",
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
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
		data := DB{}

		if err := rows.Scan(&data.TableName); err != nil {
			log.Printf("%#v\n", err)
			return nil, err
		}

		tableNames = append(tableNames, data.TableName)
	}

	return tableNames, nil
}

func internalServerError(w http.ResponseWriter) {

	response, _ := json.Marshal(&Response{
		"error": "internal server error",
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(response)

	return
}

func isExistsParam(r *http.Request, key string) bool {

	param := r.FormValue(key)
	if len(param) > 0 {
		return true
	}

	return false
}
