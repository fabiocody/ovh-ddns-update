package main

import log "github.com/sirupsen/logrus"

func HandleError(err error) {
	if err != nil {
		log.SetReportCaller(true)
		log.Fatal(err)
	}
}
