package rpc

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/uberfurrer/tradebot/logger"
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
		logger.Warning("could not read request")
	}
	var req Request
	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		logger.Warning("could not parse request")
	}
	var ep = strings.TrimPrefix(r.URL.Path, "/")
	logger.Infof("POST /%s", ep)
	if v, ok := s.Handlers[ep]; ok {
		resp := v.Process(req)
		if resp == nil {
			logger.Errorf("invalid request %v\n", req)
			return
		}
		if resp.Error != nil {
			logger.Infof("error processing request %s, id %s\n", resp.Error, *req.ID)
		}

		respData, _ := json.Marshal(resp)
		w.Write(respData)
		w.Header().Add("Content-Type", "application/json")
		return
	}
	logger.Infof("invalid endpoint /%s\n", ep)
	w.WriteHeader(404)
}

// Start starts listener
func (s *Server) Start(addr string, stop chan struct{}) {
	s.mux = http.NewServeMux()
	for k := range s.Handlers {
		s.mux.HandleFunc("/"+k, s.Handler)
	}
	l, _ := net.Listen("tcp", addr)
	logger.Infof("Starting server %s\n", addr)
	go http.Serve(l, s)
	defer l.Close()
	for _ = range stop {
		// close chan for exit
	}
	logger.Info("Stopping server...")

}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
