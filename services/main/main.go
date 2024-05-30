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
	rateLimiter := stdlib.NewMiddleware(limiter.New(memory.NewStore(), limiter.Rate{
		Limit:      100,          // limit to 100 requests per interval
		Identifier: "",           // identifier to group requests (e.g., by IP address)
		Period:     time.Second,  // interval to check the limit
	}))

	strictRateLimiter := stdlib.NewMiddleware(limiter.New(memory.NewStore(), limiter.Rate{
		Limit:      50,           // limit to 50 requests per interval
		Identifier: "",           // identifier to group requests (e.g., by IP address)
		Period:     time.Second,  // interval to check the limit
	}))

	srv := &http.Server{
		Addr:         hostPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      router,
	}

	InitHandlers(router, rateLimiter, strictRateLimiter)

	log.Println("Info: Starting Service on:", hostURI)
	log.Fatal(srv.ListenAndServe())
}
