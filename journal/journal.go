package journal

import (
	"fmt"
	"os"
	"robot/protocole"
)

// Ecrire ajoute une ligne lisible au journal de la partie.
func Ecrire(fichier *os.File, etat protocole.Etat, action string) error {
	if !etat.EnnemiVisible {
		_, err := fmt.Fprintf(fichier, "Tour %d : moi (%d,%d) %d PV, ennemi invisible, action %s\n",
			etat.Tour, etat.Moi.X, etat.Moi.Y, etat.Moi.Vie, action)
		return err
	}
	_, err := fmt.Fprintf(fichier, "Tour %d : moi (%d,%d) %d PV, ennemi (%d,%d) %d PV, action %s\n",
		etat.Tour, etat.Moi.X, etat.Moi.Y, etat.Moi.Vie,
		etat.Ennemi.X, etat.Ennemi.Y, etat.Ennemi.Vie, action)
	return err
}
