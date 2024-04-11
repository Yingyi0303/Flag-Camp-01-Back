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

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	router.Handle("/time", middleware.Handler(http.HandlerFunc(timeHandler))).Methods("GET")

	router.Handle("/discussion", middleware.Handler(http.HandlerFunc(postDiscussionHandler))).Methods("POST")
	router.Handle("/alldiscussions", middleware.Handler(http.HandlerFunc(getAllDiscussionsHandler))).Methods("GET")
	router.Handle("/mydiscussions", middleware.Handler(http.HandlerFunc(getMyDiscussionsHandler))).Methods("GET")
	router.Handle("/discussion", middleware.Handler(http.HandlerFunc(getDiscussionDetailHandler))).Methods("GET")
	router.Handle("/discussion", middleware.Handler(http.HandlerFunc(deleteDiscussionHandler))).Methods("DELETE")

	router.Handle("/reply", middleware.Handler(http.HandlerFunc(postReplyHandler))).Methods("POST")
	router.Handle("/myreplies", middleware.Handler(http.HandlerFunc(getMyRepliesHandler))).Methods("GET")
	router.Handle("/reply", middleware.Handler(http.HandlerFunc(deleteReplyHandler))).Methods("DELETE")
	
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
