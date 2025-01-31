package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httputil"
	"os"
)

// config is the configuration file for bolt-proxy
type config struct {
	ProxyDir     string                  `json:"proxy_dir"`
	HttpPort     string                  `json:"http_port"`
	TLSPort      string                  `json:"tls_port"`
	AdminUser    string                  `json:"admin_user"`
	LiveDir      string                  `json:"live_dir"`
	StageDir     string                  `json:"stage_dir"`
	CertDir      string                  `json:"cert_dir"`
	TlsCerts     tlsCerts                `json:"tls_certs"`
	ServiceRepos []string                `json:"service_repos"`
	Services     map[string]*serviceConf `json:"services"`
}

// tlsCerts are used for the tls server
type tlsCerts struct {
	Privkey   string `json:"privkey"`
	Fullchain string `json:"fullchain"`
}

type env map[string]string

// serviceConf is a type of application running on a port
type serviceConf struct {
	App          app    `json:"app"`
	GCloud       gcloud `json:"gcloud"`
	ReverseProxy *httputil.ReverseProxy
}

type app struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	DomainName string `json:"domain_name"`
	Version    string `json:"version"`
	Env        env    `json:"env"`
	Port       string `json:"port"`
	AlertsOn   bool   `json:"alertsOn"`
	TLSEnabled bool   `json:"tls_enabled"`
	Repo       string `json:"repo"`
}

type gcloud struct {
	Command   string `json:"command"`
	Zone      string `json:"zone"`
	Project   string `json:"project"`
	User      string `json:"user"`
	LiveDir   string `json:"livedir"`
	ProxyConf string `json:"proxyConf"`
}

var (
	globalHalt context.CancelFunc
	fullchain  string
	privkey    string
	httpPort   string
	tlsPort    string
	confPath   string = os.Getenv("proxConfPath")
	pc         config = config{}
)

func getHits() {
	b, err := os.ReadFile("logject.json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(b, &hitCounterByIP)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(len(hitCounterByIP), " hits in logject.json")
}

// init sets flags that tell log to log the date and line number. Init also
// reads the configuration file
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pc.Services = make(map[string]*serviceConf)
	// if len(os.Args) > 1 {
	// 	if os.Args[1] == "rebolt" {
	// 		rebolt()
	// 	}
	// }
	proxyConf()
	scan()
	getHits()
	fullchain = pc.CertDir + pc.TlsCerts.Fullchain
	privkey = pc.CertDir + pc.TlsCerts.Privkey
	httpPort = pc.HttpPort
	tlsPort = pc.TLSPort
}

func scan() {
	dir, err := os.ReadDir(pc.LiveDir)
	if err != nil {
		log.Println(err)
	}
	for _, d := range dir {
		b, err := os.ReadFile(pc.LiveDir + d.Name() + "/bolt.conf.json")
		if err != nil {
			log.Println(err)
		}
		sc := serviceConf{}
		err = json.Unmarshal(b, &sc)
		if err != nil {
			log.Println(err)
		}

		pc.Services[sc.App.DomainName] = makeProxy(&sc)
		pc.Services["www."+sc.App.DomainName] = pc.Services[sc.App.DomainName]
	}
	// startServices()
}

func proxyConf() {
	if len(confPath) < 1 {
		confPath = "/home/john/bp/prox.conf"
	}
	file, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &pc)
	if err != nil {
		log.Println(err)
	}
}
