package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackytck/lenslocked/controllers"
	"github.com/jackytck/lenslocked/middleware"
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

	// router
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	// controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)
	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{User: userMw}

	// routes
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/faq", faq).Methods("GET")

	// routes: Images
	imageHandler := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", imageHandler))

	// routes: Gallery
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.Upload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	// /galleries/:id/images/:filename/delete
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	fmt.Println("Starting the server at http://127.0.0.1:3000...")
	http.ListenAndServe(":3000", userMw.Apply(r))
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
