package main

import (
	"flag"
	"fmt"
	"os"
	"robot/carte"
	"robot/ia"
	"robot/protocole"
)

func main() {
	chemin := flag.String("carte", "", "fichier .map")
	flag.Int("joueur", 1, "numéro du joueur")
	flag.Parse()
	if *chemin == "" {
		fmt.Fprintln(os.Stderr, "indiquez --carte fichier.map")
		os.Exit(1)
	}

	c, err := carte.Charger(*chemin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lecteur := protocole.NouveauLecteur(os.Stdin)
	for {
		lignes, ok := lecteur.Suivant()
		if !ok {
			break
		}
		etat, err := protocole.Parser(lignes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("ATTENDS")
			continue
		}
		fmt.Println(ia.Decider(c, etat))
	}
	if err := lecteur.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
