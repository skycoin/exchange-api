package rpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

// Server serve all request for all exchanges
type Server struct {
	Handlers map[string]PackageHandler
	mux      *http.ServeMux
}

// Handler implenments http.Handler interface
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("could not read request")
	}
	var req Request
	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		log.Println("could not parse request")
	}
	var ep = strings.TrimPrefix(r.URL.Path, "/")
	log.Println(ep)

	if v, ok := s.Handlers[ep]; ok {
		resp := v.Process(req)
		log.Println("request executed")
		respData, err := json.Marshal(resp)
		if err != nil {
			log.Println("json marshal error", err)
		}
		log.Println(string(respData))
		w.Write(respData)
		w.Header().Add("Content-Type", "application/json")
		return
	}
	log.Println("could not find exchange", ep)
	w.WriteHeader(404)
}

// Start starts listener
func (s *Server) Start() {
	s.mux = http.NewServeMux()
	for k := range s.Handlers {
		s.mux.HandleFunc("/"+k, s.Handler)
	}
	l, _ := net.Listen("tcp", "localhost:12345")
	log.Println("Starting server on localhost:12345")
	http.Serve(l, s)

}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
