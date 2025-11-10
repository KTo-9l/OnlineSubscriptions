package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func summaryPriceForPeriodHandler(rw http.ResponseWriter, rq *http.Request) {
	log.Println("[Sum] New request")

	if rq.Method != http.MethodGet {
		log.Println("[Sum] Wrong method")
		http.Error(rw, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	startDateString := rq.URL.Query().Get("startDate")
	endDateString := rq.URL.Query().Get("endDate")
	userId := rq.URL.Query().Get("userId")
	serviceName := rq.URL.Query().Get("serviceName")

	startDate, err := parseTime(startDateString)
	if err != nil {
		log.Printf("[Sum] Error parsing startDate: %s", err.Error())
		http.Error(rw, "Error startDate", http.StatusBadRequest)
		return
	}

	endDate, err := parseTime(endDateString)
	if err != nil {
		log.Printf("[Sum] Error parsing endDate: %s", err.Error())
		http.Error(rw, "Error endDate", http.StatusBadRequest)
		return
	}

	sum, err := calculatePrice(startDate, endDate, userId, serviceName)
	if err != nil {
		log.Printf("[Sum] Error calculating price: %s", err.Error())
		http.Error(rw, "Some error while calculating", http.StatusInternalServerError)
		return
	}

	rw.Write([]byte(strconv.Itoa(int(sum))))
}

func calculatePrice(startPeriod, endPeriod time.Time, userId, serviceName string) (int, error) {
	query := `
	SELECT price, start_date, end_date
	FROM subscription
	WHERE user_id = $1 AND service_name = $2
	AND (start_date < $3)
	AND (end_date > $4 OR end_date IS NULL) 
	`

	var (
		price     int
		startDate time.Time
		endDate   *time.Time
	)
	err := db.QueryRow(query, userId, serviceName, endPeriod, startPeriod).Scan(&price, &startDate, &endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		log.Printf("[DB][Select] Error select/scan subscriptions for period with constraints: %s", err.Error())
		return 0, err
	}

	lowBorder := maxTime(startDate, startPeriod)
	highBorder := endPeriod
	if endDate != nil {
		highBorder = minTime(endPeriod, *endDate)
	}

	return ((highBorder.Year()-lowBorder.Year())*12 + (int(highBorder.Month()) - int(lowBorder.Month()))) * price, err
}
