package solver

import "push-swap/internal/stack"

// bfsNode memorise, pour un etat atteint, d'ou l'on vient et par quelle
// instruction. Cela permet de reconstruire le chemin a la fin.
type bfsNode struct {
	parent string // cle de l'etat precedent
	op     string // instruction appliquee pour arriver ici
}

// bfs cherche la sequence d'instructions la plus courte (recherche en largeur).
//
// Principe :
//  1. On part de l'etat initial et on le met dans une file (queue).
//  2. On retire un etat de la file et on lui applique chacune des 11
//     instructions : chaque resultat est un nouvel etat.
//  3. Si un nouvel etat n'a jamais ete vu, on le memorise et on l'ajoute
//     a la file. S'il est trie (a croissante, b vide), on a fini.
//
// Comme on explore d'abord tous les etats a 1 instruction, puis a 2, puis
// a 3..., la premiere solution trouvee a forcement le moins d'instructions.
func bfs(ranks []int) []string {
	start := stack.New(ranks)
	startKey := key(start)

	// visited evite de reexplorer un etat deja vu (sinon boucle infinie,
	// par exemple "ra" puis "rra" revient au meme etat).
	visited := map[string]bfsNode{startKey: {}}
	queue := []*stack.Stacks{start}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		curKey := key(cur)

		for _, op := range stack.Ops {
			next := cur.Clone()
			_ = next.Apply(op)
			nextKey := key(next)

			if _, seen := visited[nextKey]; seen {
				continue
			}
			visited[nextKey] = bfsNode{parent: curKey, op: op}

			if next.IsSolved() {
				return buildPath(visited, startKey, nextKey)
			}
			queue = append(queue, next)
		}
	}
	return nil // impossible en pratique : tout etat peut etre trie
}

// buildPath remonte de l'etat final jusqu'a l'etat initial grace aux parents,
// puis inverse la liste pour l'avoir dans l'ordre d'execution.
func buildPath(visited map[string]bfsNode, startKey, endKey string) []string {
	var ops []string
	for k := endKey; k != startKey; k = visited[k].parent {
		ops = append(ops, visited[k].op)
	}
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

// key transforme un etat en chaine unique pour pouvoir l'utiliser comme cle
// de map. Les rangs sont petits (< 6) donc chacun tient dans un octet ;
// 0xFF sert de separateur entre la pile a et la pile b.
func key(s *stack.Stacks) string {
	b := make([]byte, 0, len(s.A)+len(s.B)+1)
	for _, v := range s.A {
		b = append(b, byte(v))
	}
	b = append(b, 0xFF)
	for _, v := range s.B {
		b = append(b, byte(v))
	}
	return string(b)
}
