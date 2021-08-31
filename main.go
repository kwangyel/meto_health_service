package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
)

const (
	STATE_SERVICE = "STATE_SERVICE"
	AMQP_SERIVICE = "AMQP_SERVICE"
	DB_SERVICE    = "DATABASE_SERVICE"
	API_SERVICE   = "API_SERVICE"
)

type healthLogger struct{}

// func (l healthLogger)OnResultsUpdated(ll map[string]gosundheit.Result)
// {
// 	log.Println("There are %d results, general health is %t\n", len(results), allHealthy(results))
// }
func (l healthLogger) OnResultsUpdated(results map[string]gosundheit.Result) {
	// log.Printf("There are %d results, general health is %v\n", len(results), results)
	// log.Println(results[STATE_SERVICE].Error)
	var errorServices []string
	for key, element := range results {
		if element.Error != nil {
			log.Println("There is an error in ", element)
			errorServices = append(errorServices, key)
		}
	}
	if len(errorServices) > 0 {
		go sendMail(&EmailDTO{fromName: "root", fromEmail: "root@meto-transport.com", toName: []string{"kwangyel"}, toEmails: []string{"kwangyel@gmail.com"}, subject: "System Failure notice", body: strings.Join(errorServices, ",")})
	}
}

func main() {

	var CheckHealthLogger healthLogger

	// h := gosundheit.New()
	h := gosundheit.New(gosundheit.WithHealthListeners(CheckHealthLogger))

	httpCheckConf := checks.HTTPCheckConfig{
		CheckName: STATE_SERVICE,
		Timeout:   2 * time.Minute,
		// dependency you're checking - use your own URL here...
		// this URL will fail 50% of the times
		URL: "http://localhost:9090/state-service",
	}

	httpCheck, err := checks.NewHTTPCheck(httpCheckConf)
	if err != nil {
		fmt.Println(err)
		return // your call...
	}

	err = h.RegisterCheck(
		httpCheck,
		gosundheit.InitialDelay(time.Second),      // the check will run once after 1 sec
		gosundheit.ExecutionPeriod(2*time.Second), // the check will be executed every 10 sec
	)
	if err != nil {
		fmt.Println("Failed to register check: ", err)
		return // or whatever
	}

	http.Handle("/status", healthhttp.HandleHealthJSON(h))
	log.Println("Starting health service at port 8888")
	http.ListenAndServe(":8888", nil)
}
