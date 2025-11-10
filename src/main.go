package main

import (
	"log"
	"net/http"
)

func main() {
	loadConfig()
	initLog()

	if err := initDB(); err != nil {
		log.Panic("Error initializing database: ", err)
	}

	log.Println("Succesfully started service")

	// handler for Create, Update, List
	http.HandleFunc("/subscriptions", func(rw http.ResponseWriter, rq *http.Request) {
		switch rq.Method {
		case http.MethodPost:
			createHandler(rw, rq)
		case http.MethodPut:
			updateHandler(rw, rq)
		case http.MethodGet:
			listHandler(rw, rq)
		default:
			http.Error(rw, "Wrong method", http.StatusMethodNotAllowed)
		}
	})

	// handler for Get, Delete
	http.HandleFunc("/subscriptions/", func(rw http.ResponseWriter, rq *http.Request) {
		switch rq.Method {
		case http.MethodGet:
			getHandler(rw, rq)
		case http.MethodDelete:
			deleteHandler(rw, rq)
		default:
			http.Error(rw, "Wrong method", http.StatusMethodNotAllowed)
		}
	})

	// handler for sum for period
	http.HandleFunc("/subscriptions/sum", summaryPriceForPeriodHandler)

	http.ListenAndServe(cfg.Server.Port, nil)
}
