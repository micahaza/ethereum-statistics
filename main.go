package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type jResponse struct {
	Page int64 `json:"page,omitempty"`
	Data *Data `json:"data,omitempty"`
}

// Data ...
type Data struct {
	First  int64 `json:"first,omitempty"`
	Second int64 `json:"second,omitempty"`
	Third  int64 `json:"third,omitempty"`
}

var tpl *template.Template
var responseStruct []jResponse

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
	log.Println("INDEX...")
}

func about(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "about.gohtml", nil)
	log.Println("ABOUT...")
}

func chart(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "chart.gohtml", nil)
	log.Println("CHART...")
}

func updateChartData(w http.ResponseWriter, req *http.Request) {
	data := append(responseStruct, jResponse{Page: 2, Data: &Data{First: 23, Second: 82, Third: 22}})
	jData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	log.Println("API CALL...")
}

func main() {
	port := flag.String("p", "8080", "port")
	dir := flag.String("d", "./templates/", "dir")
	flag.Parse()
	router := mux.NewRouter()

	// Serving static files, JS, CSS, Images
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	router.HandleFunc("/", index)
	router.HandleFunc("/about", about)
	router.HandleFunc("/chart", chart)

	apiRoutes := router.PathPrefix("/api").Subrouter()
	apiRoutes.Path("/update-chart-data").Methods("GET").HandlerFunc(updateChartData)

	log.Printf("Serving %s on port %s\n", *dir, *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
