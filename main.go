package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackytck/lenslocked/controllers"
	"github.com/jackytck/lenslocked/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jacky"
	password = "natnat"
	dbname   = "lenslocked_dev"
)

func main() {
	// db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	must(err)
	defer services.Close()
	services.AutoMigrate()
	// services.DestructiveReset()

	// controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	// router
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	// routes
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/faq", faq).Methods("GET")

	// routes: Gallery
	r.Handle("/galleries/new", galleriesC.New).Methods("GET")
	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Sorry, but we couldn't find the page you were looking for!</h1>")
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Frequently Asked Questions</h1><p>Here is a list of questions that our users commonly ask.</p>")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
