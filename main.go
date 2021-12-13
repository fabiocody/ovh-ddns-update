package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	externalip "github.com/glendc/go-external-ip"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	arg.MustParse(&Args)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
	SetupDB()
	currentIP := GetCurrentIP()
	previousIP := GetPreviousIP()
	if currentIP != previousIP {
		log.Infof("IPs don't match (current is %s, previous is %s)", currentIP, previousIP)
		UpdateDDNS(currentIP)
		SaveCurrentIP(currentIP)
	} else {
		log.Infof("IPs match (%s, %s). Skipping update", currentIP, previousIP)
	}
}

func GetCurrentIP() string {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	HandleError(err)
	log.Infof("Current IP is %s", ip)
	return ip.String()
}

func UpdateDDNS(ip string) {
	url := fmt.Sprintf("https://%s:%s@www.ovh.com/nic/update?system=dyndns&hostname=%s&myip=%s", Args.OvhId, Args.OvhPassword, Args.Domain, ip)
	_, err := http.Get(url)
	HandleError(err)
	log.Infof("IP updated (%s)", ip)
}
