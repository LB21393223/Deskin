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

// / Charger lit le fichier de carte et vérifie qu'il est valide.
func Charger(chemin string) (Carte, error) {
	contenu, err := os.ReadFile(chemin)
	if err != nil {
		return Carte{}, err
	}
	nom := nomFichier(chemin)
	texte := strings.ReplaceAll(string(contenu), "\r\n", "\n")
	lignes := strings.Split(texte, "\n")
	for len(lignes) > 0 && lignes[len(lignes)-1] == "" {
		lignes = lignes[:len(lignes)-1]
	}
	if len(lignes) == 0 {
		return Carte{}, fmt.Errorf("%s : fichier vide", nom)
	}
	c, debut, err := lireEntete(lignes, nom)
	if err != nil {
		return Carte{}, err
	}
	if err := verifierEntete(c, nom); err != nil {
		return Carte{}, err
	}
	c.Grille = lignes[debut:]
	if err := verifierGrille(&c, debut, nom); err != nil {
		return Carte{}, err
	}
	return c, nil
}

// nomFichier récupère le nom du fichier sans son chemin.
func nomFichier(chemin string) string {
	parties := strings.Split(chemin, "/")
	nom := parties[len(parties)-1]
	parties = strings.Split(nom, "\\")
	return parties[len(parties)-1]
}

// lireEntete lit les informations avant la grille.
func lireEntete(lignes []string, nom string) (Carte, int, error) {
	c := Carte{}
	for i, ligne := range lignes {
		if strings.HasPrefix(ligne, "#") {
			return c, i, nil
		}
		champs := strings.Fields(ligne)
		if len(champs) == 0 {
			continue
		}
		if err := lireChamp(&c, champs, ligne, nom, i+1); err != nil {
			return Carte{}, 0, err
		}
	}
	return c, len(lignes), nil
}

// lireChamp choisit comment lire une ligne de l'en-tête.
func lireChamp(c *Carte, champs []string, ligne, nom string, n int) error {
	switch champs[0] {
	case "NOM":
		return lireNom(c, ligne, nom, n)
	case "DIFFICULTE":
		return lireDifficulte(c, champs, nom, n)
	case "TOURS":
		return lireTours(c, champs, nom, n)
	case "TAILLE":
		return lireTaille(c, champs, nom, n)
	default:
		return fmt.Errorf("%s ligne %d : clé inconnue %q", nom, n, champs[0])
	}
}

// lireNom lit le nom de la carte.
func lireNom(c *Carte, ligne, nom string, n int) error {
	if c.Nom != "" {
		return fmt.Errorf("%s ligne %d : NOM défini deux fois", nom, n)
	}
	lignePropre := strings.TrimSpace(ligne)
	c.Nom = strings.TrimSpace(lignePropre[3:])
	if c.Nom == "" {
		return fmt.Errorf("%s ligne %d : NOM attend un nom", nom, n)
	}
	return nil
}

// lireDifficulte lit et vérifie la difficulté.
func lireDifficulte(c *Carte, champs []string, nom string, n int) error {
	if c.Difficulte != 0 {
		return fmt.Errorf("%s ligne %d : DIFFICULTE défini deux fois", nom, n)
	}
	if len(champs) != 2 {
		return fmt.Errorf("%s ligne %d : DIFFICULTE attend un entier entre 1 et 3", nom, n)
	}
	valeur, err := strconv.Atoi(champs[1])
	if err != nil || valeur < 1 || valeur > 3 {
		return fmt.Errorf("%s ligne %d : DIFFICULTE attend un entier entre 1 et 3", nom, n)
	}
	c.Difficulte = valeur
	return nil
}

// lireTours lit et vérifie le nombre de tours.
func lireTours(c *Carte, champs []string, nom string, n int) error {
	if c.Tours != 0 {
		return fmt.Errorf("%s ligne %d : TOURS défini deux fois", nom, n)
	}
	if len(champs) != 2 {
		return fmt.Errorf("%s ligne %d : TOURS attend un entier positif", nom, n)
	}
	valeur, err := strconv.Atoi(champs[1])
	if err != nil || valeur <= 0 {
		return fmt.Errorf("%s ligne %d : TOURS attend un entier positif", nom, n)
	}
	c.Tours = valeur
	return nil
}

// lireTaille lit et vérifie la largeur et la hauteur.
func lireTaille(c *Carte, champs []string, nom string, n int) error {
	if c.Largeur != 0 {
		return fmt.Errorf("%s ligne %d : TAILLE défini deux fois", nom, n)
	}
	if len(champs) != 3 {
		return fmt.Errorf("%s ligne %d : TAILLE attend deux entiers", nom, n)
	}
	largeur, e1 := strconv.Atoi(champs[1])
	hauteur, e2 := strconv.Atoi(champs[2])
	if e1 != nil || e2 != nil {
		return fmt.Errorf("%s ligne %d : TAILLE attend deux entiers", nom, n)
	}
	if largeur < 3 || hauteur < 3 {
		return fmt.Errorf("%s ligne %d : TAILLE trop petite (3 x 3 minimum)", nom, n)
	}
	c.Largeur = largeur
	c.Hauteur = hauteur
	return nil
}

// verifierEntete vérifie que toutes les informations obligatoires sont présentes.
func verifierEntete(c Carte, nom string) error {
	if c.Nom == "" {
		return fmt.Errorf("%s : en-tête incomplet : NOM manquant", nom)
	}
	if c.Difficulte == 0 {
		return fmt.Errorf("%s : en-tête incomplet : DIFFICULTE manquant", nom)
	}
	if c.Tours == 0 {
		return fmt.Errorf("%s : en-tête incomplet : TOURS manquant", nom)
	}
	if c.Largeur == 0 {
		return fmt.Errorf("%s : en-tête incomplet : TAILLE manquant", nom)
	}
	return nil
}

// verifierGrille vérifie les lignes et les départs des deux joueurs.
func verifierGrille(c *Carte, debut int, nom string) error {
	if len(c.Grille) != c.Hauteur {
		return fmt.Errorf("%s : %d lignes de grille, %d attendues", nom, len(c.Grille), c.Hauteur)
	}
	depart1, depart2 := 0, 0
	for y, ligne := range c.Grille {
		if err := verifierLigne(*c, ligne, y, debut, nom, &depart1, &depart2); err != nil {
			return err
		}
	}
	if depart1 == 0 {
		return fmt.Errorf("%s : aucun départ pour le joueur 1", nom)
	}
	if depart2 == 0 {
		return fmt.Errorf("%s : aucun départ pour le joueur 2", nom)
	}
	return nil
}

// verifierLigne vérifie la largeur et les cases d'une ligne.
func verifierLigne(c Carte, ligne string, y, debut int, nom string, depart1, depart2 *int) error {
	nombreCases := 0
	for range ligne {
		nombreCases++
	}
	if nombreCases != c.Largeur {
		return fmt.Errorf("%s ligne %d : %d cases, %d attendues", nom, debut+y+1, nombreCases, c.Largeur)
	}
	colonne := 0
	for _, symbole := range ligne {
		colonne++
		if err := verifierCase(c, symbole, y, colonne, debut, nom); err != nil {
			return err
		}
		if symbole == '1' {
			*depart1 = *depart1 + 1
			if *depart1 > 1 {
				return fmt.Errorf("%s ligne %d colonne %d : départ du joueur 1 déjà défini", nom, debut+y+1, colonne)
			}
		}
		if symbole == '2' {
			*depart2 = *depart2 + 1
			if *depart2 > 1 {
				return fmt.Errorf("%s ligne %d colonne %d : départ du joueur 2 déjà défini", nom, debut+y+1, colonne)
			}
		}
	}
	return nil
}

// verifierCase vérifie le symbole et la bordure d'une case.
func verifierCase(c Carte, symbole rune, y, colonne, debut int, nom string) error {
	if !symboleValide(symbole) {
		return fmt.Errorf("%s ligne %d colonne %d : symbole inconnu %q", nom, debut+y+1, colonne, string(symbole))
	}
	if (y == 0 || y == c.Hauteur-1 || colonne == 1 || colonne == c.Largeur) && symbole != '#' {
		return fmt.Errorf("%s ligne %d colonne %d : la bordure doit être un mur", nom, debut+y+1, colonne)
	}
	return nil
}

// symboleValide indique si un symbole est autorisé dans la grille.
func symboleValide(symbole rune) bool {
	switch symbole {
	case '#', '.', '~', '^', '+', '*', '!', '$', '1', '2':
		return true
	}
	return false
}
