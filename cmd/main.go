package main

import (
	"log"

	"github.com/dmdhrumilmistry/gohtools/dns"
)

func main() {

	// TCP Servers
	// StartEchoServer(4444, "ioCopyEcho")
	// StartTcpProxy(4444, "localhost:8000")
	// StartListener(4444, HandleNcConn) // connect with server using net.Dial or nc IP 4444

	// HTTP request
	// ip := http.GetHostMachinePublicIP() // package to import: "github.com/dmdhrumilmistry/gohtools/http"
	// fmt.Printf("[*] Host Machine Public Ip: %s\n", ip)

	// msf login
	// msfHost := os.Getenv("MSF_HOST")
	// msfUsername := os.Getenv("MSF_USERNAME")
	// msfPassword := os.Getenv("MSF_PASSWORD")

	// if msfHost == "" || msfPassword == "" || msfUsername == "" {
	// 	log.Fatalln("[ERROR] Required environment variables (MSF_HOST, MSF_USERNAME, MSF_PASSWORD) not found!")
	// }

	// msf, err := metasploit.NewMsfClient(msfHost, msfUsername, msfPassword)
	// if err != nil {
	// 	log.Fatalf("[ERROR] Unable to create MSF client due to error: %s", err)
	// }
	// defer msf.Logout()

	// sessions, err := msf.ListSessions()
	// if err != nil {
	// 	log.Panicf("[ERROR] Unable to list MSF sessions due to error: %s", err)
	// }

	// if len(sessions) < 1 {
	// 	fmt.Println("[*] No active sessions found!")
	// } else {
	// 	fmt.Println("Sessions:")
	// 	for _, session := range sessions {
	// 		fmt.Printf("%5d\t%s\n", session.Id, session.Info)
	// 	}
	// }

	// Bing Doc Search
	// var docLinks []string
	// site := "example.com"
	// docType := "xlsx"

	// bsClient := bing.NewBingSearch()
	// bsClient.SearchDocument(site, docType)

	// for _, result := range bsClient.Results {
	// 	docLink, ok := result.(string)
	// 	if !ok {
	// 		continue
	// 	}
	// 	docLinks = append(docLinks, docLink)
	// }

	// fmt.Printf("[*] Docs found: %d\n", len(docLinks))
	// bsClient.ExtractDocumentsMetaData(docLinks)

	// Start Malicious Servers
	// http.HandleFunc("/hello", ghttp.EchoParam)
	// http.ListenAndServe("0.0.0.0:8000", nil) // use default server mux

	// var router ghttp.Router
	// http.ListenAndServe("0.0.0.0:8000", &router) // use custom router

	// HTTP Middleware
	// f := http.HandlerFunc(ghttp.EchoParam)
	// l := ghttp.Logger{Handler: f}
	// http.ListenAndServe("0.0.0.0:8000", &l)

	// Middleware: Negroni and Handler: gorillaMux
	// r := mux.NewRouter()
	// n := negroni.Classic()
	// n.UseHandler(r)
	// n.Use(&ghttp.TrivialMiddleware{})

	// r.HandleFunc("/hello", ghttp.EchoParam).Methods("GET")
	// r.HandleFunc("/users/{user:[a-z]+}", ghttp.EchoUserParam).Methods("GET")

	// http.ListenAndServe("0.0.0.0:8000", n)

	// Start Creds harvesting server
	// filePath := "creds.txt"
	// port := 8000
	// captureUri := "/login"
	// serverRootPath := "public"
	// credsHarvestingServ := http.NewLoginHarvest(filePath, port, captureUri, serverRootPath)
	// credsHarvestingServ.StartHarvesting()

	// KeyLogger
	// listenAddr := "localhost:8080"
	// wsAddr := listenAddr
	// jsTemplateFilePath := "logger.js"
	// keyLogger := http.NewKeyLogger(listenAddr, wsAddr, jsTemplateFilePath)

	// keyLogger.StartLogger()

	// Malicious C2C server with Reverse Proxy/Virtual Hosts
	// proxies := make(map[string]string)
	// proxies["1.gohtools.com"] = "http://localhost:8001"
	// proxies["2.gohtools.com"] = "http://localhost:8002"

	// reverseProxy := http.NewReverseProxy("0.0.0.0:80", proxies)
	// reverseProxy.StartReverseProxy()

	// DNS Server
	dnsServer := "1.1.1.1"
	domain := "microsoft.com"
	wordlist := "/Users/apple/Downloads/web-subdomains.txt"
	// wordlist := "/Users/apple/Downloads/test.txt"

	dnsFuzzer := dns.NewDnsFuzzer(dnsServer, 100, wordlist)
	results, err := dnsFuzzer.StartSubDomainFuzzer(domain)
	if err != nil {
		log.Printf("[!] Error occurred while fuzzing sub-domains: %s", err)
	}
	dnsFuzzer.PrintResults(results)

}
