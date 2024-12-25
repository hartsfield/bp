package main

import (
	"fmt"
	"log"
	"os/exec"
)

var certbotCmd string = "sudo certbot certonly --noninteractive --agree-tos " +
	"--cert-name boltorg -d hrtsfld.xyz -d walboard.xyz -d tagmachine.xy" +
	"z -d btstrmr.xyz -d bolt-marketing.org -d statui.hrtsfld.xyz -d mys" +
	"terygift.org -m johnathanhartsfield@gmail.com --standalone"

var copyCerts string = "sudo cp /etc/letsencrypt/live/boltorg/privkey.pem /etc/letsencrypt/live/boltorg/fullchain.pem tlsCerts/"

var forwardPort80to8080 string = "sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080"
var forwardPort443to8443 string = "sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 8443"

// func rebolt() {
// 	fmt.Println("Root Privileges Required...")

// 	// forward ports
// 	fmt.Println(localCommand(strings.Split(forwardPort443to8443, " ")))
// 	fmt.Println(localCommand(strings.Split(forwardPort80to8080, " ")))

// 	// get TSL certs
// 	if len(os.Args) > 2 {
// 		if os.Args[2] == "recert" {
// 			fmt.Println(localCommand(strings.Split(certbotCmd, " ")))
// 			fmt.Println(localCommand(strings.Split(copyCerts, " ")))
// 		}
// 	}
// 	// startServices()
// }

func localCommand(command []string, release bool, cdir string) string {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = cdir
	o, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("local command error: ", err, string(o))
	}
	fmt.Println("test")
	if release {
		err = cmd.Process.Release()
		if err != nil {
			log.Println(err)
		}
	}
	// fmt.Println(cmd.String(), string(o))
	return string(o)
}

// func startServices() {
// 	for domain := range pc.Services {
// 		go func() {
// 			fmt.Println(domain)
// 			if !strings.Contains(domain, "www") {
// 				live := os.Getenv("HOME") + "/live/"
// 				os.Setenv("PWD", live+domain)
// 				fmt.Println(localCommand(strings.Split("go build -o "+live+domain+"/"+domain, " ")))
// 				fmt.Println(localCommand([]string{live + domain + "/./" + domain}))
// 			}
// 		}()
// 	}
// }
