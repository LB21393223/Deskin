package tests

import (
	"robot/protocole"
	"testing"
)

// TestParserBrouillard vérifie un tour sans ennemi visible en protocole v3.
func TestParserBrouillard(t *testing.T) {
	lignes := []string{"TOUR 30", "MOI 5 5 VIE 70 MUN 3 BOUCLIER 0", "OBJET 6 7 MUNITION"}
	etat, err := protocole.Parser(lignes)
	if err != nil || etat.EnnemiVisible || len(etat.Objets) != 1 {
		t.Errorf("état = %+v, erreur = %v", etat, err)
	}
}

// TestParserV1 vérifie la lecture d'un tour en protocole v1.
func TestParserV1(t *testing.T) {
	lignes := []string{"TOUR 12", "MOI 3 4 VIE 80", "ENNEMI 7 4 VIE 40"}
	etat, err := protocole.Parser(lignes)
	if err != nil || !etat.EnnemiVisible || etat.Moi.Vie != 80 {
		t.Errorf("état = %+v, erreur = %v", etat, err)
	}
}
