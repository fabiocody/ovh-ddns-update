package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	externalip "github.com/glendc/go-external-ip"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
	"os"
)

type ArgsType struct {
	Database    string `default:".ovh-ddns-update.db" help:"database file"`
	Domain      string `arg:"positional,required" help:"OVH DynHost domain"`
	OvhId       string `arg:"positional,required" help:"OVH DynHost identifier"`
	OvhPassword string `arg:"positional,required" help:"OVH DynHost identifier password"`
}

var Args ArgsType

func main() {
	arg.MustParse(&Args)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
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

func GetPreviousIP() string {
	data, err := os.ReadFile(Args.Database)
	_, ok := err.(*fs.PathError)
	if ok {
		log.Warn("Previous IP not found")
		return "not_found"
	} else {
		HandleError(err)
	}
	return string(data)
}

func UpdateDDNS(ip string) {
	url := fmt.Sprintf("https://%s:%s@www.ovh.com/nic/update?system=dyndns&hostname=%s&myip=%s", Args.OvhId, Args.OvhPassword, Args.Domain, ip)
	_, err := http.Get(url)
	HandleError(err)
	log.Infof("IP updated (%s)", ip)
}

func SaveCurrentIP(ip string) {
	err := os.WriteFile(Args.Database, []byte(ip), 0644)
	HandleError(err)
}

func HandleError(err error) {
	if err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
	}
}
