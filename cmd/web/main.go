package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/mealies/bookings/pkg/config"
	"github.com/mealies/bookings/pkg/handlers"
	"github.com/mealies/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const localHost = "localHost:"
const portNumber = "8080"

var app config.AppConfig
var session *scs.SessionManager

// Main is the main app function
func main() {

	// change this to true when in prod
	app.InProduction = false

	// set up session using alexedwars/scs package
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // set true in prod

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on %s%s", localHost, portNumber))

	srv := &http.Server{
		Addr:    localHost + portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
