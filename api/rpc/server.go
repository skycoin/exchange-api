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
	Handlers map[string]Wrapper
	mux      *http.ServeMux
}

// Handler implenments http.Handler interface
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var req Request
	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		w.WriteHeader(500)
	}
	ep := strings.TrimPrefix(r.URL.Path, "/")
	if v, ok := s.Handlers[ep]; ok {
		resp := v.Process(req)
		if resp == nil {
			return
		}
		if resp.Error != nil {
			log.Printf("error: %s\n", resp.Error)
		}

		respData, _ := json.Marshal(resp)
		w.Write(respData)
		w.Header().Add("Content-Type", "application/json")
		return
	}
	log.Printf("not found %s\n", ep)
	w.WriteHeader(404)
}

// Start starts listener
func (s *Server) Start(addr string, stop chan struct{}) {
	s.mux = http.NewServeMux()
	for k := range s.Handlers {
		s.mux.HandleFunc("/"+k, s.Handler)
	}
	l, _ := net.Listen("tcp", addr)
	log.Printf("Starting server %s\n", addr)
	go http.Serve(l, s)
	defer l.Close()
	<-stop
	log.Println("Server stopped")

}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
