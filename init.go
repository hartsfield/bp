package main

import (
	"context"
	"encoding/json"
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
	logPath    string = os.Getenv("proxLogPath")
	pc         config = config{}
)

// init sets flags that tell log to log the date and line number. Init also
// reads the configuration file
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pc.Services = make(map[string]*serviceConf)
	fullchain = pc.CertDir + pc.TlsCerts.Fullchain
	privkey = pc.CertDir + pc.TlsCerts.Privkey
	httpPort = pc.HttpPort
	tlsPort = pc.TLSPort
}

func scan() error {
	dir, err := os.ReadDir(pc.LiveDir)
	if err != nil {

		return err
	}
	for _, d := range dir {
		b, err := os.ReadFile(pc.LiveDir + d.Name() + "/bolt.conf.json")
		if err != nil {
			// log.Println(err)
		}
		sc := serviceConf{}
		err = json.Unmarshal(b, &sc)
		if err != nil {
			return err
		}

		pc.Services[sc.App.DomainName] = makeProxy(&sc)
		pc.Services["www."+sc.App.DomainName] = pc.Services[sc.App.DomainName]
	}
	// startServices()
	return nil
}

func proxyConf() {
	if len(confPath) < 1 {
		confPath = os.Getenv("HOME") + "/bp/prox.conf"
	}
	file, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal("EXITING \nNo prox.conf found, set proxConfPath={path to prox.conf}", err)
	}
	err = json.Unmarshal(file, &pc)
	if err != nil {
		log.Println(err)
	}

	if logPath == "" {
		logPath = "./log.txt"
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("EXITING \nError opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}
