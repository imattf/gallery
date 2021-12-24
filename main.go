package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"

	"gitlab.com/go-courses/lenslocked.com/controllers"
	"gitlab.com/go-courses/lenslocked.com/email"
	"gitlab.com/go-courses/lenslocked.com/middleware"
	"gitlab.com/go-courses/lenslocked.com/models"
	"gitlab.com/go-courses/lenslocked.com/rand"
)

func main() {

	envPtr := flag.Bool("prod", false, "Provide this flag in prod. Ensures a .config file used for start-up")
	flag.Parse()

	cfg := LoadConfig(*envPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
		models.WithOAuth(),
	)

	must(err)

	defer services.Close()

	// Gorm services to reset datastore...
	// Destroy and Reset the entire datastore
	// services.DestructiveReset()
	// Auto construct from the gorm data model
	services.AutoMigrate()

	// email mailgun stuff...
	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("Gallery Support", "support@"+mgCfg.Domain),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
		// works w/ new version of mailgun
		// email.WithMailgun(mgCfg.Domain, mgCfg.APIKey),
	)

	// instance a gorilla mux
	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	configs := make(map[string]*oauth2.Config)
	configs[models.OAuthDropbox] = &oauth2.Config{
		ClientID:     cfg.Dropbox.ID,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.Dropbox.AuthURL,
			TokenURL: cfg.Dropbox.TokenURL,
		},
		RedirectURL: "http://localhost:3000/oauth/dropbox/callback",
	}
	oauthsC := controllers.NewOAuths(services.OAuth, configs)

	b, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))
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
	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	r.Handle("/forgot", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")
	// r.HandleFunc("/cookie", usersC.CookieTest).Methods("GET")

	// OAuth routes
	r.HandleFunc("/oauth/{service:[a-z]+}/connect", requireUserMw.ApplyFn(oauthsC.Connect))
	r.HandleFunc("/oauth/{service:[a-z]+}/callback", requireUserMw.ApplyFn(oauthsC.Callback))
	r.HandleFunc("/oauth/{service:[a-z]+}/test", requireUserMw.ApplyFn(oauthsC.DropboxTest))

	// Asset routes
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
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

	// Server startup...
	fmt.Printf("Starting galleries on port :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		fmt.Println("Database not found!!...")
		panic(err)
	}
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// have to do this first because Fprint's return StatusOk by default
	w.WriteHeader(http.StatusNotFound)

	// print my path
	//fmt.Fprintf(w, r.URL.Path)

	fmt.Fprint(w, "<h1>We could not find the page you are looking for :( </h1> <p>Please emaul us at <a href=\"mailto:matthew@faulkners.io\">support@gallery.faulkners.io</a> if you keep getting sent to an invalid page.</p>")

	// logging to console
	fmt.Println("404 page")
}
