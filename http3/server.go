package http3

import (
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

type Server struct {
	Addr     string
	CertPath string
	KeyPath  string
}

func NewServer(addr string, certPath string, keyPath string) *Server {
	return &Server{
		Addr:     addr,
		CertPath: certPath,
		KeyPath:  keyPath,
	}
}

func (s *Server) SetupHandler(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir("/tmp/")))

}

func (s *Server) Listen() {
	mux := http.NewServeMux()

	s.SetupHandler(mux)

	log.Printf("Starting QUIC server on %s with Cert %s and Key %s\n", s.Addr, s.CertPath, s.KeyPath)
	if err := http3.ListenAndServeQUIC(s.Addr, s.CertPath, s.KeyPath, mux); err != nil {
		log.Fatalf("failed to start QUIC server on %s with Cert %s and Key %s due to error %v\n", s.Addr, s.CertPath, s.KeyPath, err)
	}

}
