package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	addr  = ":8000"
	store *graph.Graph
)

func init() {
	store = graph.New()

	flag.StringVar(&addr, "addr", addr, "Address and port to service and accept client connections.")
	dumpfile := flag.String("dump", "", "Load a dump (.draft) file.")
	flag.Parse()

	if dumpfile != nil && *dumpfile != "" {
		log.Printf("Loading from %s", *dumpfile)
		data, err := ioutil.ReadFile(*dumpfile)
		if err != nil {
			log.Fatal(err)
		}

		dump := pb.DumpResp{}
		if err := proto.Unmarshal(data, &dump); err != nil {
			log.Fatal(err)
		}

		if err := load(store, dump); err != nil {
			log.Fatal(err)
		}
	}
}

// load a dump into the graph.
func load(g *graph.Graph, dump pb.DumpResp) error {
	for _, node := range dump.Nodes {
		if _, err := g.AddNode(node.Uid, node.Label, convertServicePropsToGraphKVs(node.Properties)...); err != nil {
			return fmt.Errorf("[load] %s", err)
		}
	}

	for _, edge := range dump.Edges {
		if _, err := g.AddEdge(edge.Uid, edge.SourceUid, edge.Label, edge.TargetUid, convertServicePropsToGraphKVs(edge.Properties)...); err != nil {
			return fmt.Errorf("[load] %s", err)
		}
	}

	return nil
}

// run start the RPC service.
func run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("[run] %s", err)
	}

	s := grpc.NewServer()
	server := &server{graph: store}
	pb.RegisterGraphServer(s, server)

	// c := make(chan os.Signal, 1)

	// go func() {
	// 	<-c
	// 	var b bytes.Buffer
	// 	server.Save(&b)
	// 	ioutil.WriteFile("../web/example/dump.draft", b.Bytes(), 0644)
	// }()

	// signal.Notify(c, os.Interrupt)
	log.Printf("[%s] Service accepting connections on %s", "run", listener.Addr())
	return s.Serve(listener)
}

// main is the main entrypoint.
func main() {
	log.Fatal(run(addr))
}
