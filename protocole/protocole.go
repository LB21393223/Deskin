package protocole

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Robot struct {
	X, Y, Vie, Mun, Bouclier int
}

type Objet struct {
	X, Y int
	Type string
}

type Etat struct {
	Tour          int
	Moi           Robot
	Ennemi        Robot
	EnnemiVisible bool
	Objets        []Objet
}

// LireBloc lit les lignes d'un tour jusqu'à FIN.
func LireBloc(scanner *bufio.Scanner) ([]string, bool) {
	lignes := []string{}
	for scanner.Scan() {
		ligne := strings.TrimSpace(scanner.Text())
		if ligne == "FIN" {
			return lignes, true
		}
		lignes = append(lignes, ligne)
	}
	return nil, false
}

// Parser transforme les lignes d'un tour en valeurs Go.
func Parser(lignes []string) (Etat, error) {
	etat := Etat{}
	vuTour, vuMoi, vuEnnemi := false, false, false

	for _, ligne := range lignes {
		champs := strings.Fields(ligne)
		if len(champs) == 0 {
			continue
		}
		switch champs[0] {
		case "TOUR":
			if vuTour || len(champs) < 2 {
				return Etat{}, fmt.Errorf("ligne TOUR invalide")
			}
			n, err := strconv.Atoi(champs[1])
			if err != nil || n < 1 {
				return Etat{}, fmt.Errorf("ligne TOUR invalide")
			}
			etat.Tour = n
			vuTour = true
		case "MOI":
			if vuMoi || len(champs) < 5 || champs[3] != "VIE" {
				return Etat{}, fmt.Errorf("ligne MOI invalide")
			}
			x, e1 := strconv.Atoi(champs[1])
			y, e2 := strconv.Atoi(champs[2])
			vie, e3 := strconv.Atoi(champs[4])
			if e1 != nil || e2 != nil || e3 != nil || vie < 0 {
				return Etat{}, fmt.Errorf("ligne MOI invalide")
			}
			mun := 10
			bouclier := 0
			if len(champs) > 5 && champs[5] == "MUN" {
				if len(champs) < 9 || champs[7] != "BOUCLIER" {
					return Etat{}, fmt.Errorf("ligne MOI invalide")
				}
				var e4, e5 error
				mun, e4 = strconv.Atoi(champs[6])
				bouclier, e5 = strconv.Atoi(champs[8])
				if e4 != nil || e5 != nil || mun < 0 || bouclier < 0 {
					return Etat{}, fmt.Errorf("ligne MOI invalide")
				}
			}
			etat.Moi = Robot{x, y, vie, mun, bouclier}
			vuMoi = true
		case "ENNEMI":
			if vuEnnemi || len(champs) < 5 || champs[3] != "VIE" {
				return Etat{}, fmt.Errorf("ligne ENNEMI invalide")
			}
			x, e1 := strconv.Atoi(champs[1])
			y, e2 := strconv.Atoi(champs[2])
			vie, e3 := strconv.Atoi(champs[4])
			if e1 != nil || e2 != nil || e3 != nil || vie < 0 {
				return Etat{}, fmt.Errorf("ligne ENNEMI invalide")
			}
			etat.Ennemi = Robot{X: x, Y: y, Vie: vie}
			etat.EnnemiVisible = true
			vuEnnemi = true
		case "OBJET":
			if len(champs) < 4 {
				return Etat{}, fmt.Errorf("ligne OBJET invalide")
			}
			x, e1 := strconv.Atoi(champs[1])
			y, e2 := strconv.Atoi(champs[2])
			if e1 != nil || e2 != nil {
				return Etat{}, fmt.Errorf("ligne OBJET invalide")
			}
			switch champs[3] {
			case "SOIN", "MUNITION", "BOUCLIER", "TRESOR":
				etat.Objets = append(etat.Objets, Objet{x, y, champs[3]})
			default:
				return Etat{}, fmt.Errorf("ligne OBJET invalide")
			}
		}
	}
	if !vuTour || !vuMoi {
		return Etat{}, fmt.Errorf("bloc incomplet")
	}
	return etat, nil
}
