package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jehiah/legislator/legistar"
	"github.com/julienschmidt/httprouter"
)

type App struct {
	legistar  *legistar.Client
	templates *template.Template
}

func (a *App) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	people, err := a.legistar.Persons()
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Unknown Error", 500)
		return
	}
	err = a.templates.ExecuteTemplate(w, "index.html", struct {
		People legistar.Persons
	}{
		People: people.Active(),
	})
	if err != nil {
		log.Printf("%s", err)
	}
}

func (a *App) Person(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("person_id"))
}

func main() {
	listen := flag.String("address", "0.0.0.0:7002", "address to listen on")
	templatePath := flag.String("templates", "templates", "path to templates")
	flag.Parse()
	app := &App{
		templates: compileTemplates(*templatePath),
		legistar: &legistar.Client{
			Client: "nyc",
			Token:  os.Getenv("NYC_LEGISLATOR_TOKEN"),
		},
	}
	router := httprouter.New()
	router.GET("/", app.Index)
	router.GET("/people/:person_id", app.Person)
	log.Printf("listening on %s", *listen)
	http.ListenAndServe(*listen, router)
}
