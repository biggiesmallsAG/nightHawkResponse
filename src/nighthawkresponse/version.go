package nighthawk

import "fmt"

var Version string = "1.0.4"
var Build string = "2786.652745"

func ShowVersion(app string) {
	fmt.Printf("\t nighthawk Response - %s\n", app)
	fmt.Printf("\t Version: %s, Build: %s\n", Version, Build)
	fmt.Println("\t By Team nighthawk (0xredskull, biggiesmalls)")
}
