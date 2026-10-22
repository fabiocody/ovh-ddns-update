package main

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	externalip "github.com/glendc/go-external-ip"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	HandleError(consensus.UseIPProtocol(4))
	ip, err := consensus.ExternalIP()
	HandleError(err)
	log.Infof("Current IP is %s", ip)
	return ip.String()
}

func GetPreviousIP() string {
	data, err := os.ReadFile(Args.Database)
	if errors.Is(err, fs.ErrNotExist) {
		log.Warn("Previous IP not found")
		return "not_found"
	} else {
		HandleError(err)
	}
	return string(data)
}

func UpdateDDNS(ip string) {
	query := url.Values{"system": {"dyndns"}, "hostname": {Args.Domain}, "myip": {ip}}
	req, err := http.NewRequest(http.MethodGet, "https://www.ovh.com/nic/update?"+query.Encode(), nil)
	HandleError(err)
	req.SetBasicAuth(Args.OvhId, Args.OvhPassword)
	res, err := http.DefaultClient.Do(req)
	HandleError(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	HandleError(err)
	response := strings.TrimSpace(string(body))
	if res.StatusCode != http.StatusOK || !(strings.HasPrefix(response, "good") || strings.HasPrefix(response, "nochg")) {
		HandleError(fmt.Errorf("OVH update failed (HTTP %d): %s", res.StatusCode, response))
	}
	log.Infof("IP updated (%s): %s", ip, response)
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
