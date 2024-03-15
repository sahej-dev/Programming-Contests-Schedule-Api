package main

import (
	"net"
	"net/http"

	"snow.sahej.io/loggers"
	ratelimiter "snow.sahej.io/rate_limiter"
)

func (app *Application) rateLimitIfRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverError(w, err)
			return
		}

		limiter := ratelimiter.GetInstance().GetLimiterForIp(ip)

		if !limiter.Allow() {
			app.tooManyRequestsError(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggers.GetInstance().InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer app.handleErrorByClosingConnection(w)
		next.ServeHTTP(w, r)
	})
}
