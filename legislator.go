package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/handlers"
	"github.com/jehiah/legislator/legistar"
	"github.com/julienschmidt/httprouter"
)

type App struct {
	legistar  *legistar.Client
	templates *template.Template

	people       legistar.Persons
	peopleBySlug map[string]legistar.Person
	voteTypes    legistar.VoteTypes
}

func NewApp(client *legistar.Client, t *template.Template) *App {
	people, err := client.Persons()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d people", len(people))
	sort.Slice(people, func(i, j int) bool { return people[i].FullName < people[j].FullName })
	voteTypes, err := client.VoteTypes()
	log.Printf("loaded %d vote types", len(voteTypes))
	if err != nil {
		panic(err)
	}

	return &App{
		legistar:     client,
		templates:    t,
		people:       people,
		peopleBySlug: people.Active().Lookup(),
		voteTypes:    voteTypes,
	}
}

func (a *App) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := a.templates.ExecuteTemplate(w, "index.html", struct {
		People legistar.Persons
	}{
		People: a.people.Active(),
	})
	if err != nil {
		log.Printf("%s", err)
	}
}

func (a *App) Person(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, ok := a.peopleBySlug[ps.ByName("person_slug")]
	if !ok {
		http.Error(w, "Person not found", 400)
		return
	}

	officeRecords, err := a.legistar.PersonOfficeRecords(p.ID)
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Unknown Error", 500)
		return
	}

	err = a.templates.ExecuteTemplate(w, "person.html", struct {
		Person        legistar.Person
		OfficeRecords legistar.OfficeRecords
	}{
		Person:        p,
		OfficeRecords: officeRecords,
	})
	if err != nil {
		log.Printf("%s", err)
	}
}

func (a *App) PersonVotes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, ok := a.peopleBySlug[ps.ByName("person_slug")]
	if !ok {
		http.Error(w, "Person not found", 400)
		return
	}

	v, err := a.legistar.PersonVotes(p.ID)
	if err != nil {
		if legistar.IsNotFoundError(err) {
			http.Error(w, "Not Found", 404)
			return
		}
		log.Printf("%s", err)
		http.Error(w, "Unknown Error", 500)
		return
	}

	err = a.templates.ExecuteTemplate(w, "person_votes.html", struct {
		Person legistar.Person
		Votes  legistar.Votes
	}{
		Person: p,
		Votes:  v,
	})
	if err != nil {
		log.Printf("%s", err)
	}
}

func main() {
	listen := flag.String("address", "0.0.0.0:7002", "address to listen on")
	templatePath := flag.String("templates", "templates", "path to templates")
	flag.Parse()
	app := NewApp(&legistar.Client{
		Client: "nyc",
		Token:  os.Getenv("NYC_LEGISLATOR_TOKEN"),
	}, compileTemplates(*templatePath))
	router := httprouter.New()
	router.GET("/", app.Index)
	router.GET("/people/:person_slug", app.Person)
	router.GET("/people/:person_slug/votes", app.PersonVotes)
	log.Printf("listening on %s", *listen)
	http.ListenAndServe(*listen, handlers.LoggingHandler(os.Stdout, router))
}
