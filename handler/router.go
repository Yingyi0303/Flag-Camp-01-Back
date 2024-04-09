package handler

import (
	"net/http"

	jwtMiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var signingKey []byte

func InitRouter() http.Handler {
	signingKey = []byte("secret")

	middleware := jwtMiddleware.New(jwtMiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()

	router.Handle("/time", middleware.Handler(http.HandlerFunc(timeHandler))).Methods("GET")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}