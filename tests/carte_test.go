package tests

import (
	"robot/carte"
	"testing"
)

type casCarte struct {
	fichier string
	message string
}

var cartesInvalides = []casCarte{
	{"02_symbole_inconnu.map", "02_symbole_inconnu.map ligne 8 colonne 5 : symbole inconnu \"?\""},
	{"05_bordure_ouverte.map", "05_bordure_ouverte.map ligne 8 colonne 1 : la bordure doit être un mur"},
	{"07_depart_double.map", "07_depart_double.map ligne 8 colonne 4 : départ du joueur 1 déjà défini"},
	{"04_taille_petite.map", "04_taille_petite.map ligne 4 : TAILLE trop petite (3 x 3 minimum)"},
}

func TestCartesInvalides(t *testing.T) {
	for _, c := range cartesInvalides {
		t.Run(c.fichier, func(t *testing.T) {
			_, err := carte.Charger("cartes/" + c.fichier)
			if err == nil || err.Error() != c.message {
				t.Errorf("erreur = %v, attendu %q", err, c.message)
			}
		})
	}
}
