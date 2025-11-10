package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func createHandler(rw http.ResponseWriter, rq *http.Request) {
	log.Println("[Create] New request")

	bodyBytes, err := io.ReadAll(rq.Body)
	if err != nil {
		log.Printf("[Create] Can't read request body: %s\n", err.Error())
		http.Error(rw, "Reading body error", http.StatusBadRequest)
		return
	}

	subscription := &Subscription{}
	err = json.Unmarshal(bodyBytes, subscription)
	if err != nil {
		log.Printf("[Create] Error unmarshalling rq.Body to struct: %s\n", err.Error())
		http.Error(rw, "Unmarshal error", http.StatusBadRequest)
		return
	}

	err = subscription.insertInDatabase()
	if err != nil {
		log.Printf("[Create] Error inserting subscription in db: %s", err.Error())
		http.Error(rw, "Error inserting in db", http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(subscription)
	if err != nil {
		log.Printf("[Create] Error marshal subscription: %s", err.Error())
		http.Error(rw, "Error marshal subscription", http.StatusInternalServerError)
		return
	}

	rw.Write(respBytes)
}

func (subscription *Subscription) insertInDatabase() (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[DB][Insert] Error begin transaction: %s", err.Error())
		return
	}

	var startDate, endDate time.Time
	startDate, err = parseTime(subscription.StartDate)
	if err != nil {
		log.Printf("[DB][Insert] Error parsing start_date: %s", err.Error())
		return
	}
	if subscription.EndDate != nil {
		endDate, err = parseTime(*subscription.EndDate)
		if err != nil {
			log.Printf("[DB][Insert] Error parsing end_date: %s", err.Error())
			return
		}
	}

	if endDate.IsZero() {
		err = tx.QueryRow("INSERT INTO subscription (service_name, price, user_id, start_date) VALUES ($1, $2, $3, $4) RETURNING id",
			subscription.ServiceName, subscription.Price, subscription.UserId, startDate).Scan(&subscription.Id)
	} else {
		err = tx.QueryRow("INSERT INTO subscription (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			subscription.ServiceName, subscription.Price, subscription.UserId, startDate, endDate).Scan(&subscription.Id)
	}

	if err != nil {
		log.Printf("[DB][Insert] Error executing a query: %s", err.Error())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[DB][Insert] Error commiting transaction: %s", err.Error())
		return
	}

	log.Printf("[DB][Insert] Successfully inserted subscription. Its id is %d", subscription.Id)
	return
}

func getHandler(rw http.ResponseWriter, rq *http.Request) {
	log.Println("[Get] New request")

	parts := strings.Split(rq.URL.Path, "/")
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		log.Printf("[Get] Error parsing requested subscription id: %s", err.Error())
		http.Error(rw, "Error parsing id", http.StatusBadRequest)
		return
	}

	subscription, err := getSubscription(id)
	if err != nil {
		log.Printf("[Get] Error get subscription by id: %s", err.Error())
		http.Error(rw, "Error get subscription by id", http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(subscription)
	if err != nil {
		log.Printf("[Get] Error marshal subscription: %s", err.Error())
		http.Error(rw, "Error marshal subscription", http.StatusInternalServerError)
		return
	}

	rw.Write(respBytes)
}

func getSubscription(id int) (subscription *Subscription, err error) {
	subscription = &Subscription{}
	err = db.QueryRow("SELECT id, service_name, price, user_id, start_date, end_date, is_deleted FROM subscription WHERE id = $1", id).Scan(
		&subscription.Id, &subscription.ServiceName, &subscription.Price, &subscription.UserId, &subscription.StartDate, &subscription.EndDate, &subscription.Deleted)

	return subscription, err
}

func updateHandler(rw http.ResponseWriter, rq *http.Request) {
	log.Println("[Update] New request")

	bodyBytes, err := io.ReadAll(rq.Body)
	if err != nil {
		log.Printf("[Update] Can't read request body: %s\n", err.Error())
		http.Error(rw, "Reading body error", http.StatusBadRequest)
		return
	}

	subscription := &Subscription{}
	err = json.Unmarshal(bodyBytes, subscription)
	if err != nil {
		log.Printf("[Update] Error unmarshalling rq.Body to struct: %s\n", err.Error())
		http.Error(rw, "Unmarshal error", http.StatusBadRequest)
		return
	}

	err = subscription.updateInDatabase()
	if err != nil {
		log.Printf("[Update] Error updating subscription in db: %s", err.Error())
		http.Error(rw, "Error updating in db", http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(subscription)
	if err != nil {
		log.Printf("[Update] Error marshal subscription: %s", err.Error())
		http.Error(rw, "Error marshal subscription", http.StatusInternalServerError)
		return
	}

	rw.Write(respBytes)
}

func (subscription *Subscription) updateInDatabase() (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[DB][Update] Error begin transaction: %s", err.Error())
		return
	}

	var startDate, endDate time.Time
	startDate, err = parseTime(subscription.StartDate)
	if err != nil {
		log.Printf("[DB][Update] Error parsing start_date: %s", err.Error())
		return
	}
	if subscription.EndDate != nil {
		endDate, err = parseTime(*subscription.EndDate)
		if err != nil {
			log.Printf("[DB][Update] Error parsing end_date: %s", err.Error())
			return
		}
	}

	if endDate.IsZero() {
		_, err = tx.Exec("UPDATE subscription SET service_name = $1, price = $2, user_id = $3, start_date = $4 WHERE id = $5",
			subscription.ServiceName, subscription.Price, subscription.UserId, startDate, subscription.Id)
	} else {
		_, err = tx.Exec("UPDATE subscription SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6",
			subscription.ServiceName, subscription.Price, subscription.UserId, startDate, endDate, subscription.Id)
	}

	if err != nil {
		log.Printf("[DB][Update] Error executing a query: %s", err.Error())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[DB][Update] Error commiting transaction: %s", err.Error())
		return
	}

	log.Println("[DB][Update] Successfully updated subscription")
	return
}

func deleteHandler(rw http.ResponseWriter, rq *http.Request) {
	log.Println("[Delete] New request")

	parts := strings.Split(rq.URL.Path, "/")
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		log.Printf("[Delete] Error parsing requested subscription id: %s", err.Error())
		http.Error(rw, "Error parsing id", http.StatusBadRequest)
		return
	}

	err = deleteSubscription(id)
	if err != nil {
		log.Printf("[Delete] Error delete subscription: %s", err.Error())
		http.Error(rw, "Error delete subscription", http.StatusInternalServerError)
		return
	}
}

func deleteSubscription(id int) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[DB][Delete] Error begin transaction: %s", err.Error())
		return
	}

	_, err = tx.Exec("UPDATE subscription SET is_deleted = true WHERE id = $1", id)
	if err != nil {
		log.Printf("[DB][Delete] Error executing a query: %s", err.Error())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[DB][Delete] Error commiting transaction: %s", err.Error())
		return
	}

	log.Println("[DB][Delete] Successfully deleted subscription")

	return
}

func listHandler(rw http.ResponseWriter, rq *http.Request) {
	log.Println("[List] New request")

	subscriptions, err := listSubscriptions()
	if err != nil {
		log.Printf("[List] Error list subscriptions: %s", err.Error())
		http.Error(rw, "Error list subscriptions", http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(subscriptions)
	if err != nil {
		log.Printf("[List] Error marshal subscriptions: %s", err.Error())
		http.Error(rw, "Error marshal subscriptions", http.StatusInternalServerError)
		return
	}

	rw.Write(respBytes)
}

func listSubscriptions() (subscriptions []Subscription, err error) {
	rows, err := db.Query("SELECT id, service_name, price, user_id, start_date, end_date FROM subscription WHERE is_deleted = false")
	if err != nil {
		log.Printf("[DB][Select] Error select subscriptions list: %s", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		subscription := Subscription{}
		err = rows.Scan(&subscription.Id, &subscription.ServiceName, &subscription.Price, &subscription.UserId, &subscription.StartDate, &subscription.EndDate)
		if err != nil {
			log.Printf("[DB][Select] Error scan row to struct: %s", err.Error())
			return
		}

		subscriptions = append(subscriptions, subscription)
	}

	return
}
