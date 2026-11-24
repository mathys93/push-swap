# push-swap

Projet en Go composé de deux programmes :

- `push-swap` calcule et affiche une liste courte d'instructions qui trie les entiers reçus en argument.
- `checker` lit des instructions sur l'entrée standard, les exécute et affiche `OK` si la pile `a` est triée et la pile `b` vide, `KO` sinon.

Seule la bibliothèque standard de Go est utilisée.

## Compilation

Linux / macOS :

```sh
go build -o push-swap ./cmd/push-swap
go build -o checker ./cmd/checker
```

Windows (Git Bash ou PowerShell) :

```sh
go build -o push-swap.exe ./cmd/push-swap
go build -o checker.exe ./cmd/checker
```

## Utilisation

```sh
$ ./push-swap "2 1 3 6 5 8"
pb
pb
ra
sa
rrr
pa
pa

$ ARG="4 67 3 87 23"; ./push-swap "$ARG" | wc -l
6

$ ARG="4 67 3 87 23"; ./push-swap "$ARG" | ./checker "$ARG"
OK

$ ./push-swap "0 one 2 3"
Error

$ echo -e "rra\npb\nsa\nrra\npa" | ./checker "3 2 1 0"
OK
```

Les nombres peuvent être donnés dans un seul argument (`"3 2 1"`) ou dans plusieurs (`3 2 1`). Le premier nombre est le sommet de la pile `a`.

## Tests

```sh
go test ./...
go test ./... -v    # affiche aussi le nombre moyen d'instructions par taille
```

## Arborescence

```
push-swap/
  go.mod
  cmd/
    push-swap/main.go     programme push-swap : parsing -> tri -> affichage
    checker/main.go       programme checker : parsing -> lecture stdin -> OK/KO
  internal/
    parser/               lecture et validation des arguments (entiers, doublons)
    stack/                les piles a et b et les 11 instructions
    solver/
      solver.go           point d'entrée du tri, choix de la stratégie
      bfs.go              tri optimal pour 6 nombres ou moins
      greedy.go           tri glouton par coût minimal pour plus de 6 nombres
```

Chaque package a une seule responsabilité :

| Package  | Rôle                                                              |
|----------|-------------------------------------------------------------------|
| parser   | transformer les arguments en `[]int`, détecter les erreurs        |
| stack    | représenter les piles et appliquer les instructions               |
| solver   | calculer la liste d'instructions                                  |
| cmd/...  | orchestrer et afficher (stdout pour le résultat, stderr pour Error) |

`push-swap` et `checker` utilisent le même package `stack` : les instructions se comportent donc forcément pareil dans les deux programmes.

## Représentation des piles

Une pile est un slice `[]int` dont l'indice 0 est le sommet. Les instructions :

| Instruction | Effet                                                  |
|-------------|--------------------------------------------------------|
| sa / sb     | échange les deux premiers éléments de a / b            |
| ss          | sa et sb en même temps                                 |
| pa / pb     | prend le sommet de b et le met sur a / l'inverse       |
| ra / rb     | rotation : le premier élément devient le dernier       |
| rr          | ra et rb en même temps                                 |
| rra / rrb   | rotation inverse : le dernier devient le premier       |
| rrr         | rra et rrb en même temps                               |

Une instruction sur une pile vide ou trop petite ne fait rien (pas de plantage).

## Algorithme

Avant de trier, chaque nombre est remplacé par son rang (0 pour le plus petit, 1 pour le suivant...). Exemple : `42 -7 100` devient `1 0 2`. L'ordre relatif est le même, donc les instructions aussi, mais les calculs sont plus simples.

### Cas 1 : 6 nombres ou moins, recherche en largeur (BFS)

On considère chaque configuration des deux piles comme un "état". Depuis l'état de départ, on applique chacune des 11 instructions, ce qui donne de nouveaux états, puis on recommence depuis ces états, etc.

On explore d'abord tous les états atteignables en 1 instruction, puis en 2, puis en 3... La première fois qu'on tombe sur un état trié, on est donc sûr d'avoir le nombre minimal d'instructions. On garde, pour chaque état, l'état précédent et l'instruction utilisée, ce qui permet de reconstruire le chemin.

Un ensemble `visited` évite de repasser par un état déjà vu (sinon `ra` puis `rra` ferait boucler à l'infini).

Pour 6 nombres il y a environ 35 000 états possibles : c'est instantané. Au-delà, le nombre d'états explose (factorielle), d'où le second algorithme.

### Cas 2 : plus de 6 nombres, algorithme glouton par coût minimal

1. On pousse les éléments de `a` vers `b` jusqu'à ce qu'il en reste 3 dans `a`. La pile `b` est gardée triée en ordre décroissant (à une rotation près). À chaque tour, on calcule pour chaque élément de `a` le nombre d'instructions nécessaires pour :
   - l'amener au sommet de `a` ;
   - faire tourner `b` pour que l'élément arrive juste au-dessus du plus grand nombre plus petit que lui.

   On choisit l'élément le moins cher.
2. On trie les 3 éléments restants dans `a` en 2 instructions maximum.
3. On ramène les éléments de `b` vers `a` de la même manière, en choisissant toujours le moins cher à insérer au bon endroit.
4. On fait tourner `a` pour amener le plus petit nombre au sommet.

Pour calculer le coût, on compare 4 possibilités (les deux piles tournent vers le haut, les deux vers le bas, ou chacune dans un sens différent). Quand les deux piles tournent dans le même sens, on utilise `rr` ou `rrr`, qui comptent pour une seule instruction au lieu de deux.

### Résultats

Nombre moyen d'instructions mesuré par les tests (listes aléatoires) :

| Taille | Moyenne | Pire cas observé |
|--------|---------|------------------|
| 5      | optimal (BFS) | -          |
| 10     | 24      | 30               |
| 100    | 564     | 595              |
| 500    | 5119    | 5300             |

## Gestion des erreurs

`Error` suivi de `\n` est écrit sur la sortie d'erreur, avec un code de sortie 1, quand :

- un argument n'est pas un entier (`one`, `1.5`, `12abc`, nombre trop grand) ;
- un nombre apparaît deux fois ;
- (checker uniquement) une instruction n'existe pas ou est mal formatée (`SA`, `" sa"`, `"sa "`).

Sans argument, les deux programmes n'affichent rien. Si la liste est déjà triée, `push-swap` n'affiche rien (0 instruction).
