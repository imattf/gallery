package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"gitlab.com/go-courses/lenslocked.com/controllers"
	"gitlab.com/go-courses/lenslocked.com/middleware"
	"gitlab.com/go-courses/lenslocked.com/models"
	"gitlab.com/go-courses/lenslocked.com/rand"
)

// temp for dev purposes
const (
	host   = "localhost"
	port   = 5432
	user   = "matthew"
	dbname = "lenslocked_dev"
)

// var (
// 	homeView    *views.View
// 	contactView *views.View
// 	faqView     *views.View
// )

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// have to do this first because Fprint's return StatusOk by default
	w.WriteHeader(http.StatusNotFound)

	// print my path
	//fmt.Fprintf(w, r.URL.Path)

	fmt.Fprint(w, "<h1>We could not find the page you are looking for :( </h1> <p>Please emaul us at <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a> if you keep getting sent to an invalid page.</p>")

	// logging to console
	fmt.Println("404 page")
}

func main() {

	// temp for dev purposes
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	services, err := models.NewServices(psqlInfo)
	must(err)

	defer services.Close()
	// services.DestructiveReset()
	services.AutoMigrate()

	// instance a gorilla mux
	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	// TODO: Update this to be a config variable
	isProd := false
	b, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	// use custom 404 page
	r.NotFoundHandler = http.HandlerFunc(notFoundPage)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	// r.HandleFunc("/cookie", usersC.CookieTest).Methods("GET")

	//Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery routes
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	fmt.Println("Starting lenslocked on port :3000...")
	http.ListenAndServe(":3000", csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		fmt.Println("Did you start the database?")
		panic(err)
	}
}
