// Package parser transforme les arguments de la ligne de commande en une liste
// d'entiers et valide cette liste (entiers valides, pas de doublons).
//
// Responsabilite unique : lire et valider l'entree. Le parser ne connait rien
// aux piles ni aux instructions.
package parser

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidInput est renvoyee quand un argument n'est pas un entier ou quand
// un nombre apparait deux fois. Les programmes affichent alors "Error".
var ErrInvalidInput = errors.New("invalid input")

// ParseArgs convertit les arguments en entiers.
//
// Le sujet montre l'appel ./push-swap "2 1 3 6 5 8" : tous les nombres sont
// dans UN seul argument. On accepte aussi ./push-swap 2 1 3 : chaque argument
// est donc decoupe sur les espaces avec strings.Fields, qui ignore aussi les
// espaces multiples.
//
// L'ordre est conserve : le premier nombre lu sera le sommet de la pile a.
func ParseArgs(args []string) ([]int, error) {
	var values []int
	// La map sert d'ensemble : la recherche d'un doublon se fait en O(1)
	// au lieu de reparcourir toute la liste a chaque nombre.
	seen := make(map[int]bool)

	for _, arg := range args {
		for _, field := range strings.Fields(arg) {
			// strconv.Atoi refuse "one", "1.5", "12abc"... et les nombres
			// trop grands pour un int : tout cela devient une erreur.
			n, err := strconv.Atoi(field)
			if err != nil {
				return nil, ErrInvalidInput
			}
			if seen[n] {
				return nil, ErrInvalidInput
			}
			seen[n] = true
			values = append(values, n)
		}
	}
	return values, nil
}
