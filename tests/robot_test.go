package tests

import (
	"robot/carte"
	"robot/ia"
	"robot/protocole"
	"testing"
)

func TestActions(t *testing.T) {
	ouverte := carte.Carte{Largeur: 7, Hauteur: 5, Grille: []string{
		"#######", "#.....#", "#.....#", "#.....#", "#######",
	}}
	detour := carte.Carte{Largeur: 7, Hauteur: 5, Grille: []string{
		"#######", "#..#..#", "#.....#", "#..#..#", "#######",
	}}
	fermee := carte.Carte{Largeur: 7, Hauteur: 5, Grille: []string{
		"#######", "#..#..#", "#..#..#", "#..#..#", "#######",
	}}
	cas := []struct {
		nom    string
		carte  carte.Carte
		etat   protocole.Etat
		action string
	}{
		{"tir est", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 2, Mun: 1}, Ennemi: protocole.Robot{X: 3, Y: 2}}, "TIRE E"},
		{"tir ouest", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 5, Y: 2, Mun: 1}, Ennemi: protocole.Robot{X: 3, Y: 2}}, "TIRE O"},
		{"tir nord", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 3, Y: 3, Mun: 1}, Ennemi: protocole.Robot{X: 3, Y: 1}}, "TIRE N"},
		{"tir sud", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 3, Y: 1, Mun: 1}, Ennemi: protocole.Robot{X: 3, Y: 3}}, "TIRE S"},
		{"chercher munition", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 2}, Ennemi: protocole.Robot{X: 5, Y: 3}, Objets: []protocole.Objet{{X: 2, Y: 2, Type: "MUNITION"}}}, "AVANCE E"},
		{"chercher soin", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 2, Vie: 30, Mun: 1}, Ennemi: protocole.Robot{X: 5, Y: 3}, Objets: []protocole.Objet{{X: 2, Y: 2, Type: "SOIN"}}}, "AVANCE E"},
		{"tir avant soin", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 2, Vie: 30, Mun: 1}, Ennemi: protocole.Robot{X: 3, Y: 2}, Objets: []protocole.Objet{{X: 2, Y: 2, Type: "SOIN"}}}, "TIRE E"},
		{"avancer sans munition", ouverte, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 2}, Ennemi: protocole.Robot{X: 4, Y: 2}}, "AVANCE E"},
		{"contourner un mur", detour, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 1, Mun: 1}, Ennemi: protocole.Robot{X: 5, Y: 1}}, "AVANCE S"},
		{"ennemi inaccessible", fermee, protocole.Etat{EnnemiVisible: true, Moi: protocole.Robot{X: 1, Y: 1, Mun: 1}, Ennemi: protocole.Robot{X: 5, Y: 1}}, "ATTENDS"},
		{"brouillard v3", ouverte, protocole.Etat{Moi: protocole.Robot{X: 1, Y: 1, Mun: 1}}, "RADAR"},
	}
	for _, c := range cas {
		t.Run(c.nom, func(t *testing.T) {
			action := ia.Decider(c.carte, c.etat)
			if action != c.action {
				t.Errorf("action = %q, attendu %q", action, c.action)
			}
		})
	}
}

func TestParserBrouillard(t *testing.T) {
	lignes := []string{"TOUR 30", "MOI 5 5 VIE 70 MUN 3 BOUCLIER 0", "OBJET 6 7 MUNITION"}
	etat, err := protocole.Parser(lignes)
	if err != nil || etat.EnnemiVisible || len(etat.Objets) != 1 {
		t.Errorf("état = %+v, erreur = %v", etat, err)
	}
}

func TestParserV1(t *testing.T) {
	lignes := []string{"TOUR 12", "MOI 3 4 VIE 80", "ENNEMI 7 4 VIE 40"}
	etat, err := protocole.Parser(lignes)
	if err != nil || !etat.EnnemiVisible || etat.Moi.Vie != 80 {
		t.Errorf("état = %+v, erreur = %v", etat, err)
	}
}

func TestCartesInvalides(t *testing.T) {
	cas := []struct {
		fichier string
		message string
	}{
		{"02_symbole_inconnu.map", "02_symbole_inconnu.map ligne 8 colonne 5 : symbole inconnu \"?\""},
		{"05_bordure_ouverte.map", "05_bordure_ouverte.map ligne 8 colonne 1 : la bordure doit être un mur"},
		{"07_depart_double.map", "07_depart_double.map ligne 8 colonne 4 : départ du joueur 1 déjà défini"},
		{"04_taille_petite.map", "04_taille_petite.map ligne 4 : TAILLE trop petite (3 x 3 minimum)"},
	}
	for _, c := range cas {
		t.Run(c.fichier, func(t *testing.T) {
			_, err := carte.Charger("cartes/" + c.fichier)
			if err == nil || err.Error() != c.message {
				t.Errorf("erreur = %v, attendu %q", err, c.message)
			}
		})
	}
}
