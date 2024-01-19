package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

type ReverseProxy struct {
	HostProxies map[string]string
	proxies     map[string]*httputil.ReverseProxy
	ListenAddr  string
}

func NewReverseProxy(listenAddr string, hostProxies map[string]string) *ReverseProxy {
	proxies := make(map[string]*httputil.ReverseProxy)

	for host, proxyUrl := range hostProxies {
		remote, err := url.Parse(proxyUrl)
		if err != nil {
			log.Printf("[HTTP-REVERSE-PROXY-WARNING] Unable to parse proxy url: %s\n.", proxyUrl)
			continue
		}

		proxies[host] = httputil.NewSingleHostReverseProxy(remote)
	}

	return &ReverseProxy{
		HostProxies: hostProxies,
		proxies:     proxies,
		ListenAddr:  listenAddr,
	}
}

func (rp *ReverseProxy) StartReverseProxy() {
	r := mux.NewRouter()
	for host, proxy := range rp.proxies {
		r.Host(host).Handler(proxy)
	}

	log.Printf("[HTTP-REVERSE-PROXY-INFO] Starting Reverse proxy on address: %s\n", rp.ListenAddr)
	log.Fatalln(http.ListenAndServe(rp.ListenAddr, r))
}
