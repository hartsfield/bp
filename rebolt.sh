#!/bin/bash
cd 
for d in $HOME/live/*/ ; do
    cd $d
    go build -o $HOME/bin/$(basename ${d})
    $(basename ${d}) &
done
## var certbotCmd string = "sudo certbot certonly --noninteractive --agree-tos " +
## 	"--cert-name boltorg -d hrtsfld.xyz -d walboard.xyz -d tagmachine.xy" +
## 	"z -d btstrmr.xyz -d bolt-marketing.org -d statui.hrtsfld.xyz -d mys" +
## 	"terygift.org -m johnathanhartsfield@gmail.com --standalone"

## var copyCerts string = "sudo cp /etc/letsencrypt/live/boltorg/privkey.pem /etc/letsencrypt/live/boltorg/fullchain.pem tlsCerts/"

## var forwardPort80to8080 string = "sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080"
## var forwardPort443to8443 string = "sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 8443"

## func rebolt() {
## 	fmt.Println("Root Privileges Required...")

## 	// forward ports
## 	fmt.Println(localCommand(strings.Split(forwardPort443to8443, " ")))
## 	fmt.Println(localCommand(strings.Split(forwardPort80to8080, " ")))

## 	// get TSL certs
## 	if len(os.Args) > 2 {
## 		if os.Args[2] == "recert" {
## 			fmt.Println(localCommand(strings.Split(certbotCmd, " ")))
## 			fmt.Println(localCommand(strings.Split(copyCerts, " ")))
## 		}
## 	}
## 	// startServices()
## }


