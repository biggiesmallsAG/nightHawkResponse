/* nighthawk.logger.main
 * author: 0xredskull && biggiesmalls
 * Team nightHawk.
 *
 * Platform logging process.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	nhc "nighthawk/common"
	"nighthawklogger/config"
	"nighthawklogger/logger"
	"os"
)

type RuntimeOptions struct {
	Debug, Help bool
}

func fUsage() {
	fmt.Printf("\tnightHawk Logger.GO version %s (07/01/2017), by Team nightHawk.\n", config.VERSION)
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {

	flag.Usage = fUsage

	var runopt RuntimeOptions

	flag.BoolVar(&runopt.Help, "h", false, "Display use flags.")
	flag.BoolVar(&runopt.Debug, "d", false, "Turn on console level debugging.")

	flag.Parse()

	if runopt.Help {
		fUsage()
	}

	if runopt.Debug {
		config.DEBUG = true
	}

	var l logger.Logger
	ch, rconfig := logger.InitMQLogger()

	logs := logger.ConsumeMQLogger(ch, &rconfig)

	db := config.InitDB()
	config.CreateTable(db)

	for log := range logs {
		nhc.ConsoleMessage("main", "DEBUG", string(log.Body), config.DEBUG)

		json.Unmarshal(log.Body, &l)
		config.StoreLogs(db, &l)
	}
}
