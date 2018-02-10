package rpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Server represents an RPC server for exchange-api
type Server struct {
	Handlers map[string]Wrapper
	server   http.Server
}

// Start starts an exchange-api RPC Server on the given address, and runs until the given channel is closed
func (s *Server) Start(addr string, stop chan struct{}) {
	var mux = http.NewServeMux()

	for k := range s.Handlers {
		mux.HandleFunc("/"+k, s.handler)
	}

	s.server = http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server %s\n", addr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Error while running RPC server: %s\n", err)
		}
	}()

	<-stop

	log.Println("Stopping server...")
	if err := s.server.Close(); err != nil {
		log.Printf("Error while closing server: %s\n", err)
	}

	log.Println("Stopped server.")
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {

	// reading request body
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	defer func() {
		if err = r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var req Request
	if err = json.Unmarshal(reqBytes, &req); err != nil {
		w.WriteHeader(500)
		return
	}

	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	wrapper, ok := s.Handlers[endpoint]
	if !ok {
		log.Printf("not found %s\n", endpoint)
		w.WriteHeader(404)
		return
	}

	response := wrapper.Process(req)
	if response == nil {
		w.WriteHeader(500)
		return
	}

	if response.Error != nil {
		log.Printf("error: %s\n", response.Error)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err = w.Write(responseBytes); err != nil {
		log.Printf("Failed to write responseBytes: %s\n", err)
		w.WriteHeader(500)
		return
	}
}
