package tests

import (
	"robot/carte"
	"robot/ia"
	"robot/protocole"
	"testing"
)

type casAction struct {
	nom     string
	carte   carte.Carte
	etat    protocole.Etat
	attendu string
}

var ouverte = carte.Carte{Largeur: 7, Hauteur: 5, Grille: []string{
	"#######", "#.....#", "#.....#", "#.....#", "#######",
}}
var detour = carte.Carte{Largeur: 7, Hauteur: 5, Grille: []string{
	"#######", "#..#..#", "#.....#", "#..#..#", "#######",
}}
var fermee = carte.Carte{Largeur: 7, Hauteur: 5, Grille: []string{
	"#######", "#..#..#", "#..#..#", "#..#..#", "#######",
}}

var actions = []casAction{
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

func TestActions(t *testing.T) {
	for _, c := range actions {
		t.Run(c.nom, func(t *testing.T) {
			action := ia.Decider(c.carte, c.etat)
			if action != c.attendu {
				t.Errorf("action = %q, attendu %q", action, c.attendu)
			}
		})
	}
}
