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

	db = inmemory.New()

	fmt.Println("Inititalizing data")
	err := db.InitData()
	if err != nil {
		log.Fatal(err)
	}

	// API нашего сервера
	http.HandleFunc("/api/v1/securities", securitiesHandler)

	// запускаем веб сервер на порту 8080
	const addr = ":8080"
	fmt.Println("Starting web server at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}

func securitiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, SessionID")

	if r.Method != http.MethodGet {
		return
	}

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
