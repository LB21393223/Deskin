package main

import (
	"flag"
	"fmt"
	"os"
	"robot/carte"
)

func main() {
	chemin := flag.String("carte", "", "fichier .map")
	flag.Parse()

	if *chemin == "" {
		fmt.Fprintln(os.Stderr, "indiquez --carte fichier.map")
		os.Exit(1)
	}

	_, err := carte.Charger(*chemin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
