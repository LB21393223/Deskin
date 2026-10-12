# Deskin - Arena

Robot de combat en Go pour le projet Arena. Il lit une carte, reçoit l'état de chaque tour et répond avec une action.

## Installation

Depuis le dossier `Deskin` :

```powershell
go build -o robot.exe .
```

Le projet utilise seulement la bibliothèque standard de Go.

## Utilisation

Depuis le dossier `Deskin` :

```powershell
..\arbitre.exe --protocole 3 --carte ..\cartes\niveau1\arene.map --robot1 .\robot.exe --robot2 cible --vitesse 0 --silencieux
```

Pour essayer une de nos cartes, remplacer le chemin par `cartes\carre.map` ou `cartes\passage.map`.

## Organisation du code

| Dossier | Rôle |
| --- | --- |
| `carte` | Lit et vérifie le fichier `.map` |
| `protocole` | Lit les blocs de l'arbitre et les transforme en valeurs Go |
| `ia` | Choisit l'action du tour ; utilise une recherche de chemin simple (BFS) |
| `journal` | Écrit les tours joués dans `partie.log` |
| `main.go` | Relie ces parties et répond à l'arbitre |

## Décisions du robot

Le robot tire si l'ennemi est visible, à portée, sans mur entre eux, et s'il a des munitions. Sinon il cherche un soin quand sa vie est basse, ou des munitions quand il n'en a plus. Il se déplace vers l'ennemi en contournant les murs. Quand l'ennemi est caché par le brouillard du protocole v3, il utilise le radar. Cette stratégie est volontairement simple.

## Vérification

Depuis `Deskin` :

```powershell
go test ./...
```

Les tests vérifient des actions attendues, une entrée sans ennemi en v3 et quatre cartes invalides. Les fichiers `cartes/carre.map` et `cartes/passage.map` sont les deux cartes créées pour le groupe.

Le fichier `partie.log` est créé automatiquement quand le robot joue une partie.