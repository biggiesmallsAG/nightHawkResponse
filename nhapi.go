/*
	nightHawkAPI.main;
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"nighthawk"
	api "nighthawkapi/api/core"
	routes "nighthawkapi/api/routes"
	"os"
)

type RuntimeOptions struct {
	Debug, Help bool
	Server      string
	Port        int
	Version     bool
}

func fUsage() {
	fmt.Printf("\tnightHawkAPI v%s, by Team nightHawk (Daniel Eden & Roshan Maskey).\n", api.VERSION)
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {

	flag.Usage = fUsage

	var runopt RuntimeOptions

	flag.BoolVar(&runopt.Help, "h", false, "Display use flags.")
	flag.BoolVar(&runopt.Debug, "d", false, "Turn on console level debugging.")
	flag.StringVar(&runopt.Server, "s", "localhost", "Bind server to address. Default: localhost")
	flag.IntVar(&runopt.Port, "p", 8080, "Bind server to port. Default: 8080")
	flag.BoolVar(&runopt.Version, "version", false, "Show version information")

	flag.Parse()

	if runopt.Help {
		fUsage()
	}

	if runopt.Version {
		nighthawk.ShowVersion("API Server")
		os.Exit(0)
	}

	if runopt.Debug {
		api.DEBUG = true
	}

	if runopt.Server != "" || runopt.Port != 8080 {

		go api.Manager.Start()
		router := routes.NewRouter()

		api.LogDebug(api.DEBUG, fmt.Sprintf("[-] Serving on %s", fmt.Sprintf("%s:%d", runopt.Server, runopt.Port)))
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", runopt.Server, runopt.Port), router))
	}
}
