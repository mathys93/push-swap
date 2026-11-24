// Package solver calcule une liste courte d'instructions qui trie la pile a.
//
// Deux strategies selon la taille de l'entree :
//
//   - 6 nombres ou moins : recherche en largeur (BFS, voir bfs.go). On explore
//     tous les etats possibles, couche par couche, donc la premiere solution
//     trouvee est la PLUS COURTE possible (solution optimale).
//
//   - plus de 6 nombres : le BFS devient trop lent (le nombre d'etats explose),
//     on utilise un algorithme glouton "par cout" (voir greedy.go) qui donne
//     de tres bons resultats sans etre forcement optimal.
package solver

import (
	"sort"

	"push-swap/internal/stack"
)

// bfsLimit est la taille maximale traitee par le BFS. Pour 6 nombres il y a
// environ 35 000 etats a explorer, ce qui reste instantane.
const bfsLimit = 6

// Solve renvoie les instructions qui trient values (values[0] = sommet de a).
func Solve(values []int) []string {
	// Cas limite : deja triee (ou 0/1 element) -> aucune instruction.
	if stack.IsSorted(values) {
		return nil
	}

	// On remplace chaque nombre par son rang (0 = le plus petit). Exemple :
	// [42 -7 100] devient [1 0 2]. L'ordre relatif ne change pas, donc les
	// instructions trouvees sont les memes, mais on travaille sur des petits
	// entiers consecutifs, plus simples a manipuler.
	ranks := normalize(values)

	if len(ranks) <= bfsLimit {
		return bfs(ranks)
	}
	return greedy(ranks)
}

// normalize remplace chaque valeur par son rang dans la liste triee.
func normalize(values []int) []int {
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)

	rankOf := make(map[int]int, len(sorted))
	for i, v := range sorted {
		rankOf[v] = i
	}

	ranks := make([]int, len(values))
	for i, v := range values {
		ranks[i] = rankOf[v]
	}
	return ranks
}

// recorder applique des instructions sur les piles ET les enregistre.
// Le solveur simule le tri sur de vraies piles : on est sur que la liste
// enregistree produit exactement l'etat simule.
type recorder struct {
	s   *stack.Stacks
	ops []string
}

// do execute l'instruction op "times" fois et l'ajoute a la liste.
func (r *recorder) do(op string, times int) {
	for i := 0; i < times; i++ {
		// Les instructions utilisees par le solveur sont toujours valides,
		// l'erreur ne peut donc pas se produire ici.
		_ = r.s.Apply(op)
		r.ops = append(r.ops, op)
	}
}
