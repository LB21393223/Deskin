package protocole

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Robot struct {
	X        int
	Y        int
	Vie      int
	Mun      int
	Bouclier int
}

type Objet struct {
	X    int
	Y    int
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
	vu := map[string]bool{}
	for _, ligne := range lignes {
		champs := strings.Fields(ligne)
		if len(champs) == 0 {
			continue
		}
		cle := champs[0]
		if cle == "TOUR" || cle == "MOI" || cle == "ENNEMI" {
			if vu[cle] {
				return Etat{}, fmt.Errorf("ligne %s invalide", cle)
			}
			vu[cle] = true
		}
		if err := ajouter(&etat, champs); err != nil {
			return Etat{}, err
		}
	}
	if !vu["TOUR"] || !vu["MOI"] {
		return Etat{}, fmt.Errorf("bloc incomplet")
	}
	return etat, nil
}

func ajouter(etat *Etat, champs []string) error {
	switch champs[0] {
	case "TOUR":
		n, err := lireTour(champs)
		etat.Tour = n
		return err
	case "MOI":
		robot, err := lireMoi(champs)
		etat.Moi = robot
		return err
	case "ENNEMI":
		robot, err := lireEnnemi(champs)
		etat.Ennemi = robot
		etat.EnnemiVisible = err == nil
		return err
	case "OBJET":
		objet, err := lireObjet(champs)
		if err == nil {
			etat.Objets = append(etat.Objets, objet)
		}
		return err
	}
	return nil
}

func lireTour(champs []string) (int, error) {
	if len(champs) < 2 {
		return 0, fmt.Errorf("ligne TOUR invalide")
	}
	n, err := strconv.Atoi(champs[1])
	if err != nil || n < 1 {
		return 0, fmt.Errorf("ligne TOUR invalide")
	}
	return n, nil
}

func lireMoi(champs []string) (Robot, error) {
	if len(champs) < 5 || champs[3] != "VIE" {
		return Robot{}, fmt.Errorf("ligne MOI invalide")
	}
	x, e1 := strconv.Atoi(champs[1])
	y, e2 := strconv.Atoi(champs[2])
	vie, e3 := strconv.Atoi(champs[4])
	if e1 != nil || e2 != nil || e3 != nil || vie < 0 {
		return Robot{}, fmt.Errorf("ligne MOI invalide")
	}
	robot := Robot{X: x, Y: y, Vie: vie, Mun: 10}
	if len(champs) > 5 && champs[5] == "MUN" {
		if len(champs) < 9 || champs[7] != "BOUCLIER" {
			return Robot{}, fmt.Errorf("ligne MOI invalide")
		}
		robot.Mun, e1 = strconv.Atoi(champs[6])
		robot.Bouclier, e2 = strconv.Atoi(champs[8])
		if e1 != nil || e2 != nil || robot.Mun < 0 || robot.Bouclier < 0 {
			return Robot{}, fmt.Errorf("ligne MOI invalide")
		}
	}
	return robot, nil
}

func lireEnnemi(champs []string) (Robot, error) {
	if len(champs) < 5 || champs[3] != "VIE" {
		return Robot{}, fmt.Errorf("ligne ENNEMI invalide")
	}
	x, e1 := strconv.Atoi(champs[1])
	y, e2 := strconv.Atoi(champs[2])
	vie, e3 := strconv.Atoi(champs[4])
	if e1 != nil || e2 != nil || e3 != nil || vie < 0 {
		return Robot{}, fmt.Errorf("ligne ENNEMI invalide")
	}
	return Robot{X: x, Y: y, Vie: vie}, nil
}

func lireObjet(champs []string) (Objet, error) {
	if len(champs) < 4 {
		return Objet{}, fmt.Errorf("ligne OBJET invalide")
	}
	x, e1 := strconv.Atoi(champs[1])
	y, e2 := strconv.Atoi(champs[2])
	if e1 != nil || e2 != nil {
		return Objet{}, fmt.Errorf("ligne OBJET invalide")
	}
	switch champs[3] {
	case "SOIN", "MUNITION", "BOUCLIER", "TRESOR":
		return Objet{X: x, Y: y, Type: champs[3]}, nil
	}
	return Objet{}, fmt.Errorf("ligne OBJET invalide")
}
