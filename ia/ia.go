package ia

import (
	"robot/carte"
	"robot/protocole"
)

type Pos struct{ X, Y int }

func voisin(p Pos, direction string) Pos {
	switch direction {
	case "N":
		p.Y--
	case "S":
		p.Y++
	case "E":
		p.X++
	case "O":
		p.X--
	}
	return p
}

func libre(c carte.Carte, p Pos) bool {
	if p.X < 0 || p.X >= c.Largeur || p.Y < 0 || p.Y >= c.Hauteur {
		return false
	}
	return c.Grille[p.Y][p.X] != '#'
}

// PremierPas cherche le premier mouvement d'un plus court chemin.
func PremierPas(c carte.Carte, depart, cible Pos) (string, int) {
	if !libre(c, depart) || !libre(c, cible) {
		return "", -1
	}
	if depart == cible {
		return "", 0
	}
	type noeud struct {
		p       Pos
		premier string
		dist    int
	}
	vu := map[Pos]bool{depart: true}
	file := []noeud{{depart, "", 0}}
	for len(file) > 0 {
		n := file[0]
		file = file[1:]
		for _, direction := range []string{"N", "S", "E", "O"} {
			p := voisin(n.p, direction)
			if vu[p] || !libre(c, p) {
				continue
			}
			vu[p] = true
			premier := n.premier
			if premier == "" {
				premier = direction
			}
			if p == cible {
				return premier, n.dist + 1
			}
			file = append(file, noeud{p, premier, n.dist + 1})
		}
	}
	return "", -1
}

func tirer(c carte.Carte, moi, ennemi protocole.Robot) string {
	depart := Pos{moi.X, moi.Y}
	cible := Pos{ennemi.X, ennemi.Y}
	for _, direction := range []string{"N", "S", "E", "O"} {
		p := depart
		for distance := 1; distance <= 5; distance++ {
			p = voisin(p, direction)
			if !libre(c, p) {
				break
			}
			if p == cible {
				return "TIRE " + direction
			}
		}
	}
	return ""
}

func versObjet(c carte.Carte, depart Pos, objets []protocole.Objet, typeObjet string) string {
	meilleureDistance := 0
	meilleureDirection := ""
	for _, objet := range objets {
		if objet.Type != typeObjet {
			continue
		}
		direction, distance := PremierPas(c, depart, Pos{objet.X, objet.Y})
		if distance > 0 && (meilleureDistance == 0 || distance < meilleureDistance) {
			meilleureDistance = distance
			meilleureDirection = direction
		}
	}
	if meilleureDirection == "" {
		return ""
	}
	return "AVANCE " + meilleureDirection
}

// Decider choisit une action simple à partir de la carte et de l'état du tour.
func Decider(c carte.Carte, etat protocole.Etat) string {
	depart := Pos{etat.Moi.X, etat.Moi.Y}
	if etat.Moi.Mun > 0 {
		if action := tirer(c, etat.Moi, etat.Ennemi); action != "" {
			return action
		}
	}
	if etat.Moi.Vie < 50 {
		if action := versObjet(c, depart, etat.Objets, "SOIN"); action != "" {
			return action
		}
	}
	if etat.Moi.Mun == 0 {
		if action := versObjet(c, depart, etat.Objets, "MUNITION"); action != "" {
			return action
		}
		if action := versObjet(c, depart, etat.Objets, "TRESOR"); action != "" {
			return action
		}
	}
	direction, distance := PremierPas(c, depart, Pos{etat.Ennemi.X, etat.Ennemi.Y})
	if distance > 1 {
		return "AVANCE " + direction
	}
	return "ATTENDS"
}
