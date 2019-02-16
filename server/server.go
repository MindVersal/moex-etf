// Package main реализует веб-сервер проекта moex-etf
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"moex_etf/server/storage"
	"moex_etf/server/storage/inmemory"
	"net/http"
)

var db storage.Interface

func main() {

	// здесь мы можем, например, добавить проверку флагов запуска или переменной окружения
	// для выбора поставщика хранилища. выбрали память
	db = inmemory.New()

	fmt.Println("Inititalizing data")
	// инициализация данных хранилища
	err := db.InitData()
	if err != nil {
		log.Fatal(err)
	}

	// API нашего сервера
	http.HandleFunc("/api/v1/securities", securitiesHandler) // список бумаг с котировками
	http.HandleFunc("/api/v1/inflation", inflationHandler)   // инфляция по месяцам

	// запускаем веб сервер на порту 8080
	const addr = ":8080"
	fmt.Println("Starting web server at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}

// обработчик запроса котировок
func securitiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method != http.MethodGet {
		return
	}

	fmt.Println(db)

	securities, err := db.Securities()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	err = json.NewEncoder(w).Encode(securities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

}

// обработчик запроса инфляции
func inflationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method != http.MethodGet {
		return
	}

	err := json.NewEncoder(w).Encode(inflation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

}

// инфляция в России по месяцам
type inflationType struct {
	Year   int
	Values [12]float64
}

var inflation = []inflationType{
	{
		Year:   2013,
		Values: [12]float64{0.97, 0.56, 0.34, 0.51, 0.66, 0.42, 0.82, 0.14, 0.21, 0.57, 0.56, 0.51},
	},
	{
		Year:   2014,
		Values: [12]float64{0.59, 0.70, 1.02, 0.90, 0.90, 0.62, 0.49, 0.24, 0.65, 0.82, 1.28, 2.62},
	},
	{
		Year:   2015,
		Values: [12]float64{3.85, 2.22, 1.21, 0.46, 0.35, 0.19, 0.80, 0.35, 0.57, 0.74, 0.75, 0.77},
	},
	{
		Year:   2016,
		Values: [12]float64{0.96, 0.63, 0.46, 0.44, 0.41, 0.36, 0.54, 0.01, 0.17, 0.43, 0.44, 0.40},
	},
	{
		Year:   2017,
		Values: [12]float64{0.62, 0.22, 0.13, 0.33, 0.37, 0.61, 0.07, -0.54, -0.15, 0.20, 0.22, 0.42},
	},
	{
		Year:   2018,
		Values: [12]float64{0.31, 0.21, 0.29, 0.38, 0.38, 0.49, 0.27, 0.01, 0.16, 0.35, 0.50, 0.84},
	},
}
