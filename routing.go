package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var hitCounterByIP map[string]*hit = make(map[string]*hit)

type hit struct {
	IP       string
	Count    int
	Requests []*remoteReq
}
type remoteReq struct {
	Time      time.Time
	Referer   string
	Method    string
	Host      string
	URL       string
	Pattern   string
	Proto     string
	UserAgent string
	Port      string
}

// newServerConf returns a new server configuration, and is used when
// instantiating the servers that intercept HTTP/S traffic.
func newServerConf(port string, hf http.HandlerFunc) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           hf,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    0,
	}
}

// startHTTPServer is used to start the HTTP server
func startHTTPServer(s *http.Server) {
	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	// if the server encounters an error, this function will be called to
	// halt the server.
	globalHalt()
}

// startTLSServer is used to start the TLS server
func startTLSServer(s *http.Server) {
	err := s.ListenAndServeTLS(fullchain, privkey)
	if err != nil {
		log.Println(err)
	}
	// if the server encounters an error, this function will be called to
	// halt the server.
	globalHalt()
}

// forwardTLS is the handler used for requests on the secure port (TLS/HTTPS).
// forwardTLS will check if a host exists and has TLS enabled, if both are
// true, it serves the website, if the host exists, but doesn't have TLS
// enabled, it forwarss it to the HTTP server, otherwise it sends the client to
// the 'not found' page.
func forwardTLS(w http.ResponseWriter, r *http.Request) {
	hitInfo(r)
	if host, ok := pc.Services[r.Host]; ok {
		if pc.Services[r.Host].App.TLSEnabled {
			host.ReverseProxy.ServeHTTP(w, r)
			return
		}
		forwardHTTP(w, r)
		return
	}
	notFound(w, r)
}

func printLogJSON() {
	b, err := json.MarshalIndent(hitCounterByIP, "", "    ")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("logject.json", b, 0666)
	if err != nil {
		log.Println(err)
	}
	log.Println("Wrote to: logject.json")

}

func hitInfo(r *http.Request) {
	ra_ := strings.Split(r.RemoteAddr, ":")
	ra := ra_[0]
	var port_ string = ra_[1]
	log.Println("<:"+ra, r.Referer(), r.Method, "<:"+r.Host, "<:"+r.URL.String(), r.Pattern, r.Proto, "<:"+r.UserAgent())

	rr := &remoteReq{
		time.Now(),
		r.Referer(),
		r.Method,
		r.Host,
		r.URL.Host,
		r.Pattern,
		r.Proto,
		r.UserAgent(),
		port_,
	}
	if hitCounterByIP[ra] == nil {
		hitCounterByIP[ra] = &hit{
			ra,
			0,
			[]*remoteReq{},
		}
	}
	hitCounterByIP[ra].Count = hitCounterByIP[ra].Count + 1
	hitCounterByIP[ra].Requests = append(hitCounterByIP[ra].Requests, rr)
	secret := "/" + os.Getenv("secretp")
	if strings.Contains(r.URL.Path, secret) && len(secret) == len(r.URL.Path) {
		printLogJSON()
	}
}

// forwardHTTP checks the host name of HTTP traffic, if TLS is enabled, it
// re-writes the address and forwards the client to the the https website,
// other wise it forwards it to the appropriate service
func forwardHTTP(w http.ResponseWriter, r *http.Request) {
	if host, ok := pc.Services[r.Host]; ok {
		if pc.Services[r.Host].App.TLSEnabled {
			rHost := r.Host
			if r.Host[0:4] == "www." {
				rHost = rHost[4:]
			}
			target := "https://" + rHost + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
			return
		}
		hitInfo(r)
		host.ReverseProxy.ServeHTTP(w, r)
		return
	}
	notFound(w, r)
}

// notFound is used If the user tries to visit a host that can't be found.
func notFound(w http.ResponseWriter, r *http.Request) {
	hitInfo(r)
	_, err := w.Write([]byte("dreams --of=infinity && gift --of=eternity && offspring --of=UNLIMITED && TRANSCEND DESTINY %"))
	if err != nil {
		log.Println(err)
	}
}
