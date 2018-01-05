package nighthawk

import "fmt"

var Version string = "1.0.4"
var Build string = "2244.581"

func ShowVersion() {
	fmt.Println("\t nighthawk Response")
	fmt.Printf("\t Version: %s, Build: %s\n", Version, Build)
	fmt.Println("\t By Team nighthawk (0xredskull, biggiesmalls)")
}
