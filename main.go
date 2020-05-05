package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	pb "github.com/jenmud/draft/service"
	"google.golang.org/grpc"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

var (
	addr      = ":8080"
	templates = packr.New("templates", "./templates")
	static    = packr.New("static", "./static")
	client    pb.GraphClient
)

// connectRPC connects the client to the gRPC server
func connectRPC(addr string) error {
	log.Printf("Connecting to RPC server %s", addr)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {
		return err
	}

	client = pb.NewGraphClient(conn)
	log.Printf("Connected to RPC server %s", addr)
	return nil
}

func init() {
	flag.StringVar(&addr, "addr", addr, "Address and port to service and accept client connections.")
	rpcAddr := flag.String("rpc-addr", ":8000", "RPC server to connect to and communicate with.")
	flag.Parse()

	if err := connectRPC(*rpcAddr); err != nil {
		log.Fatal(err)
	}
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

	dump, err := client.Query(r.Context(), &pb.QueryReq{Query: query})
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
	dump, err := client.Dump(r.Context(), &pb.DumpReq{})
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
	dump, err := client.Stats(r.Context(), &pb.StatsReq{})
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
func run(address string) error {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/assets/json", assetJSON).Methods("GET")
	router.HandleFunc("/assets/json", assetJSONQuery).Methods("POST")
	router.HandleFunc("/stats/json", statsJSON)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static)))
	log.Printf("[%s] Service accepting connections on %s", "run", address)
	return http.ListenAndServe(address, router)
}

// main is the main entrypoint.
func main() {
	log.Fatal(run(addr))
}
