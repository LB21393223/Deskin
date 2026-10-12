package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"robot/carte"
	"robot/ia"
	"robot/journal"
	"robot/protocole"
)

func main() {
	chemin := flag.String("carte", "", "fichier .map")
	flag.Int("joueur", 1, "numéro du joueur")
	flag.Parse()
	c, err := carte.Charger(*chemin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fichier, err := journal.Ouvrir("partie.log")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	jouer(c, fichier)
	if fichier != nil {
		if err := journal.Fermer(fichier); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func jouer(c carte.Carte, fichier *os.File) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for {
		lignes, ok := protocole.LireBloc(scanner)
		if !ok {
			break
		}
		etat, err := protocole.Parser(lignes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("ATTENDS")
			continue
		}
		action := ia.Decider(c, etat)
		fmt.Println(action)
		if fichier != nil {
			if err := journal.Ecrire(fichier, etat, action); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
