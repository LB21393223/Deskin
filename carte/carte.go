package carte

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Carte struct {
	Nom        string
	Difficulte int
	Tours      int
	Largeur    int
	Hauteur    int
	Grille     []string
}

// Charger lit et vérifie un fichier de carte.
func Charger(chemin string) (Carte, error) {
	contenu, err := os.ReadFile(chemin)
	if err != nil {
		return Carte{}, err
	}

	texte := strings.ReplaceAll(string(contenu), "\r\n", "\n")
	lignes := strings.Split(texte, "\n")
	for len(lignes) > 0 && lignes[len(lignes)-1] == "" {
		lignes = lignes[:len(lignes)-1]
	}

	parties := strings.Split(chemin, "/")
	nomFichier := parties[len(parties)-1]
	parties = strings.Split(nomFichier, "\\")
	nomFichier = parties[len(parties)-1]

	if len(lignes) == 0 {
		return Carte{}, fmt.Errorf("%s : fichier vide", nomFichier)
	}

	carte := Carte{}
	debutGrille := len(lignes)

	for i, ligne := range lignes {
		if strings.HasPrefix(ligne, "#") {
			debutGrille = i
			break
		}

		champs := strings.Fields(ligne)
		if len(champs) == 0 {
			continue
		}

		switch champs[0] {
		case "NOM":
			if carte.Nom != "" {
				return Carte{}, fmt.Errorf("%s ligne %d : NOM défini deux fois", nomFichier, i+1)
			}
			lignePropre := strings.TrimSpace(ligne)
			carte.Nom = strings.TrimSpace(lignePropre[3:])
			if carte.Nom == "" {
				return Carte{}, fmt.Errorf("%s ligne %d : NOM attend un nom", nomFichier, i+1)
			}

		case "DIFFICULTE":
			if carte.Difficulte != 0 {
				return Carte{}, fmt.Errorf("%s ligne %d : DIFFICULTE défini deux fois", nomFichier, i+1)
			}
			if len(champs) != 2 {
				return Carte{}, fmt.Errorf("%s ligne %d : DIFFICULTE attend un entier entre 1 et 3", nomFichier, i+1)
			}
			valeur, err := strconv.Atoi(champs[1])
			if err != nil || valeur < 1 || valeur > 3 {
				return Carte{}, fmt.Errorf("%s ligne %d : DIFFICULTE attend un entier entre 1 et 3", nomFichier, i+1)
			}
			carte.Difficulte = valeur

		case "TOURS":
			if carte.Tours != 0 {
				return Carte{}, fmt.Errorf("%s ligne %d : TOURS défini deux fois", nomFichier, i+1)
			}
			if len(champs) != 2 {
				return Carte{}, fmt.Errorf("%s ligne %d : TOURS attend un entier positif", nomFichier, i+1)
			}
			valeur, err := strconv.Atoi(champs[1])
			if err != nil || valeur <= 0 {
				return Carte{}, fmt.Errorf("%s ligne %d : TOURS attend un entier positif", nomFichier, i+1)
			}
			carte.Tours = valeur

		case "TAILLE":
			if carte.Largeur != 0 {
				return Carte{}, fmt.Errorf("%s ligne %d : TAILLE défini deux fois", nomFichier, i+1)
			}
			if len(champs) != 3 {
				return Carte{}, fmt.Errorf("%s ligne %d : TAILLE attend deux entiers", nomFichier, i+1)
			}
			largeur, errLargeur := strconv.Atoi(champs[1])
			hauteur, errHauteur := strconv.Atoi(champs[2])
			if errLargeur != nil || errHauteur != nil {
				return Carte{}, fmt.Errorf("%s ligne %d : TAILLE attend deux entiers", nomFichier, i+1)
			}
			if largeur < 3 || hauteur < 3 {
				return Carte{}, fmt.Errorf("%s ligne %d : TAILLE trop petite (3 x 3 minimum)", nomFichier, i+1)
			}
			carte.Largeur = largeur
			carte.Hauteur = hauteur

		default:
			return Carte{}, fmt.Errorf("%s ligne %d : clé inconnue %q", nomFichier, i+1, champs[0])
		}
	}

	if carte.Nom == "" {
		return Carte{}, fmt.Errorf("%s : en-tête incomplet : NOM manquant", nomFichier)
	}
	if carte.Difficulte == 0 {
		return Carte{}, fmt.Errorf("%s : en-tête incomplet : DIFFICULTE manquant", nomFichier)
	}
	if carte.Tours == 0 {
		return Carte{}, fmt.Errorf("%s : en-tête incomplet : TOURS manquant", nomFichier)
	}
	if carte.Largeur == 0 {
		return Carte{}, fmt.Errorf("%s : en-tête incomplet : TAILLE manquant", nomFichier)
	}

	carte.Grille = lignes[debutGrille:]
	if len(carte.Grille) != carte.Hauteur {
		return Carte{}, fmt.Errorf("%s : %d lignes de grille, %d attendues", nomFichier, len(carte.Grille), carte.Hauteur)
	}

	depart1 := 0
	depart2 := 0

	for y, ligne := range carte.Grille {
		nombreCases := 0
		for range ligne {
			nombreCases++
		}
		if nombreCases != carte.Largeur {
			return Carte{}, fmt.Errorf("%s ligne %d : %d cases, %d attendues", nomFichier, debutGrille+y+1, nombreCases, carte.Largeur)
		}

		colonne := 0
		for _, symbole := range ligne {
			colonne++

			if symbole != '#' && symbole != '.' && symbole != '~' && symbole != '^' && symbole != '+' && symbole != '*' && symbole != '!' && symbole != '$' && symbole != '1' && symbole != '2' {
				return Carte{}, fmt.Errorf("%s ligne %d colonne %d : symbole inconnu %q", nomFichier, debutGrille+y+1, colonne, string(symbole))
			}
			if (y == 0 || y == carte.Hauteur-1 || colonne == 1 || colonne == carte.Largeur) && symbole != '#' {
				return Carte{}, fmt.Errorf("%s ligne %d colonne %d : la bordure doit être un mur", nomFichier, debutGrille+y+1, colonne)
			}
			if symbole == '1' {
				depart1++
				if depart1 > 1 {
					return Carte{}, fmt.Errorf("%s ligne %d colonne %d : départ du joueur 1 déjà défini", nomFichier, debutGrille+y+1, colonne)
				}
			}
			if symbole == '2' {
				depart2++
				if depart2 > 1 {
					return Carte{}, fmt.Errorf("%s ligne %d colonne %d : départ du joueur 2 déjà défini", nomFichier, debutGrille+y+1, colonne)
				}
			}
		}
	}

	if depart1 == 0 {
		return Carte{}, fmt.Errorf("%s : aucun départ pour le joueur 1", nomFichier)
	}
	if depart2 == 0 {
		return Carte{}, fmt.Errorf("%s : aucun départ pour le joueur 2", nomFichier)
	}

	return carte, nil
}
