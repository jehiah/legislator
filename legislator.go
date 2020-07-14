package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/handlers"
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

	people = people.Active()
	sort.Slice(people, func(i, j int) bool { return people[i].FullName < people[j].FullName })

	err = a.templates.ExecuteTemplate(w, "index.html", struct {
		People legistar.Persons
	}{
		People: people,
	})
	if err != nil {
		log.Printf("%s", err)
	}
}

func (a *App) Person(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("person_id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid Person ID", 400)
		return
	}

	p, err := a.legistar.Person(id)
	if err != nil {
		if legistar.IsNotFoundError(err) {
			http.Error(w, "Not Found", 404)
			return
		}
		log.Printf("%s", err)
		http.Error(w, "Unknown Error", 500)
		return
	}

	fmt.Fprintf(w, "hello, %#v!\n", p)
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
	http.ListenAndServe(*listen, handlers.LoggingHandler(os.Stdout, router))
}
