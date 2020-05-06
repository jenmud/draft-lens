package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	pb "github.com/jenmud/draft/service"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/v2"
	microConfig "github.com/micro/go-micro/v2/config"
	microEnv "github.com/micro/go-micro/v2/config/source/env"
	microFlag "github.com/micro/go-micro/v2/config/source/flag"
	microWeb "github.com/micro/go-micro/v2/web"
)

var (
	version     = microWeb.Version("v0.0.0")
	templates   = packr.New("templates", "./templates")
	static      = packr.New("static", "./static")
	config      microConfig.Config
	draftClient pb.GraphService
)

func parseArgs() {
	var err error

	config, err = microConfig.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	flag.String("addr", ":", "Address to accept client connections on")
	flag.String("name", "draft.web.srv", "Service name")
	flag.String("draft-service", "draft.srv", "Draft service name")
	flag.String("draft-addr", "", "Draft RPC server, if not given service discovery is used")
	flag.Parse()

	err = config.Load(
		microEnv.NewSource(microEnv.WithPrefix("DRAFT_LENS")),
		microFlag.NewSource(microFlag.IncludeUnset(true)),
	)

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	parseArgs()

	var service micro.Service

	srv := config.Get("draft", "addr").String("")
	if srv != "" {
		service = micro.NewService(micro.Address(srv))
	} else {
		service = micro.NewService(micro.Name(config.Get("draft", "service").String("draft.web.srv")))
	}

	draftClient = pb.NewGraphService(
		config.Get("draft", "service").String("draft.web.srv"),
		service.Client(),
	)
}

// index serves up the index page.
func index(w http.ResponseWriter, r *http.Request) {
	html, err := templates.Find("base.tmpl")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(html)
}

// assetJSONPOST serves JSON assets with a query.
func assetJSONQuery(w http.ResponseWriter, r *http.Request) {
	// parse the post form
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	query := r.FormValue("cypher")
	log.Printf("Form data %v: %s", r.Form, query)

	dump, err := draftClient.Query(r.Context(), &pb.QueryReq{Query: query})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(dump)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// assetJSON serves JSON assets.
func assetJSON(w http.ResponseWriter, r *http.Request) {
	dump, err := draftClient.Dump(r.Context(), &pb.DumpReq{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(dump)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// statsJSON serves JSON stats.
func statsJSON(w http.ResponseWriter, r *http.Request) {
	dump, err := draftClient.Stats(r.Context(), &pb.StatsReq{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(dump)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// run start the RPC service.
func run() error {
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/assets/json", assetJSON).Methods("GET")
	router.HandleFunc("/assets/json", assetJSONQuery).Methods("POST")
	router.HandleFunc("/stats/json", statsJSON)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static)))

	service := microWeb.NewService(
		microWeb.Name(config.Get("name").String("draft.web.srv")),
		version,
		microWeb.Address(config.Get("addr").String(":")),
		microWeb.Handler(router),
	)

	return service.Run()
}

// main is the main entrypoint.
func main() {
	log.Fatal(run())
}
