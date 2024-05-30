package main

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
)

func main() {
	defer logFile.Close()
	log.SetOutput(logFile)

	router := httprouter.New()

	// Create rate limiter middleware instances
	rateLimiter := stdlib.NewMiddleware(limiter.NewIPRateLimiter(memory.NewStore(), limiter.IPConfig{
		Max:        100,          // limit each IP to 100 requests per interval
		Identifier: "",           // identify clients by their IP address
		Every:      time.Second,  // interval to check the limit
	}))
	strictRateLimiter := stdlib.NewMiddleware(limiter.NewIPRateLimiter(memory.NewStore(), limiter.IPConfig{
		Max:        50,           // limit each IP to 50 requests per interval
		Identifier: "",           // identify clients by their IP address
		Every:      time.Second,  // interval to check the limit
	}))

	srv := &http.Server{
		Addr:         hostPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      router,
	}

	InitHandlers(router, rateLimiter, strictRateLimiter)

	/* https to http redirection
	go http.ListenAndServe(":80", http.HandlerFunc(httpsRedirect))
	log.Println("Info: Starting Service on:", hostURI)
	log.Fatal(srv.ListenAndServeTLS("fullchain.pem", "privkey.pem"))
	*/
	log.Println("Info: Starting Service on:", hostURI)
	log.Fatal(srv.ListenAndServe())
}
