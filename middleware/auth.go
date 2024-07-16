package middleware

import (
	"log"
	"net/http"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing fake Authorize middleware")
		urlPath := r.URL.Path
		log.Print(urlPath)
		if urlPath == "/recipe/new/" {
			log.Print("We would need to do some auth here")
		}
		next.ServeHTTP(w, r)
	})
}

//TODO Some authenticating once we set it up.
