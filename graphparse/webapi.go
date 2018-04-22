package graphparse

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"go/ast"
	"go/token"
	"strconv"
	"fmt"
	"bytes"
	"strings"

	"github.com/gorilla/mux"
)

func WebAPI(port string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(corsEnabledHeaders)

	
	src := router.PathPrefix("/src").Subrouter().Methods("GET")
	src.Path("/ast-range").
	Queries(
		"start", "{start:[0-9]+}",
		"end",   "{end:[0-9]+}",
	).HandlerFunc(getASTRange)
	src.Path("/").HandlerFunc(getSrc)


	graph := router.PathPrefix("/graph").Subrouter().Methods("GET")
	graph.Path("/last-generated").HandlerFunc(getGraphLastGenerated)

	graph.Path("/code-thread").
	Queries(
		"from", "{from:[0-9]+}",
		"to",   "{to:[0-9]+}",
	).HandlerFunc(getGraphCodeThread)

	graph.Path("/contracted").
	Queries(
		"q", "{q}",
	).HandlerFunc(getContractedGraph)

	graph.Path("/").HandlerFunc(getGraph)


	log.Fatal(http.ListenAndServe("0.0.0.0:" + port, router))
}


// Middleware
// ==========

func corsEnabledHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}



// Source
// ================

type astRangeRes struct {
	Output string `json:"output"`
}

func getASTRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	posStart, _ := strconv.ParseInt(vars["start"], 0, 64)
	posEnd, _ := strconv.ParseInt(vars["end"], 0, 64)
	fmt.Println("Showing code from", vars["start"], "to", vars["end"])

	_, nodes, _ := prog.PathEnclosingInterval(token.Pos(posStart), token.Pos(posEnd))
	buf := new(bytes.Buffer)
	ast.Fprint(buf, fset, nodes[:1], ast.NotNilFilter)
	
	res := astRangeRes{
		Output: buf.String(),
	}
	json.NewEncoder(w).Encode(res)
}

func getSrc(w http.ResponseWriter, r *http.Request) {
	// TODO for testing remove later.
	for f, src := range fileLookup {
		if f == "chain.go" {
			json.NewEncoder(w).Encode(src)
			return
		}
	}

	http.Error(w, "example file not found", http.StatusBadRequest)
}



// Graph
// ===============

func getGraphLastGenerated(w http.ResponseWriter, r *http.Request) {
	buf := bufio.NewWriter(w)
	buf.WriteString(Graph.generatedAt)
	buf.Flush()
}


type codeThreadRes struct {
	Edges []edge `json:"edges"`
}

func getGraphCodeThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fromId, _ := strconv.ParseInt(vars["from"], 0, 64)
	toId, _ := strconv.ParseInt(vars["to"], 0, 64)

	from, ok := nodeLookup[nodeid(fromId)]
	if !ok {
		http.Error(w, "from node not found", http.StatusBadRequest)
		return
	}

	to, ok := nodeLookup[nodeid(toId)]
	if !ok {
		http.Error(w, "to node not found", http.StatusBadRequest)
		return
	}

	edges := pathEnclosingNodes(Graph, from, to)
	res := codeThreadRes{
		Edges: edges,
	}
	json.NewEncoder(w).Encode(res)
}


type getContractedGraphReq struct {
	NodeTypesHidden []NodeType `json:"nodeTypesHidden"`
}

func getContractedGraph(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var req getContractedGraphReq
	if err := json.NewDecoder(strings.NewReader(vars["q"])).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nodeTypeHidden := make(map[NodeType]bool)
	for _, nodeType := range req.NodeTypesHidden {
		nodeTypeHidden[nodeType] = true
	}
	
	edges := Graph.contractEdges(func(n Node) bool {
		if nodeTypeHidden[n.Variant()] {
			return true
		}
		return false
	})

	res := Graph._toJson(edges)
	json.NewEncoder(w).Encode(res)
}


func getGraph(w http.ResponseWriter, r *http.Request) {
	res := Graph.toJson()
	json.NewEncoder(w).Encode(res)
}