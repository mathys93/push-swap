package solver

import "push-swap/internal/stack"

// greedy trie les entrees de plus de 6 nombres avec un algorithme glouton
// "par cout minimal".
//
// Etapes :
//
//  1. On pousse les elements de a vers b jusqu'a ce qu'il en reste 3 dans a.
//     La pile b est maintenue triee en ordre DECROISSANT (a une rotation pres).
//     A chaque tour, on calcule pour CHAQUE element de a combien
//     d'instructions il faudrait pour le placer au bon endroit dans b, et on
//     choisit le moins cher (d'ou le nom "glouton").
//
//  2. On trie les 3 elements restants dans a (au plus 2 instructions).
//
//  3. On ramene les elements de b vers a, toujours en choisissant celui qui
//     coute le moins cher a inserer au bon endroit dans a.
//
//  4. La pile a est alors triee mais "tournee" : on fait tourner a pour
//     amener le plus petit element au sommet.
//
// L'optimisation principale : quand a ET b doivent tourner dans le meme sens,
// on utilise "rr" ou "rrr" qui font les deux rotations en une seule
// instruction.
func greedy(ranks []int) []string {
	r := &recorder{s: stack.New(ranks)}

	// Etape 1 : a -> b
	for len(r.s.A) > 3 {
		if len(r.s.B) == 0 {
			// b vide : pas de position a chercher, on pousse directement.
			r.do("pb", 1)
			continue
		}
		m := cheapestMove(r.s.A, r.s.B, targetInDescending)
		applyMove(r, m.src, m.dst)
		r.do("pb", 1)
	}

	// Etape 2 : tri des 3 elements restants
	sortThree(r)

	// Etape 3 : b -> a (ici la source est b et la destination a)
	for len(r.s.B) > 0 {
		m := cheapestMove(r.s.B, r.s.A, targetInAscending)
		applyMove(r, m.dst, m.src)
		r.do("pa", 1)
	}

	// Etape 4 : amener le minimum au sommet par le chemin le plus court
	rotateTimes(r, "ra", "rra", shortestRotation(indexOfMin(r.s.A), len(r.s.A)))
	return r.ops
}

// move decrit les rotations necessaires avant un push.
// Convention : nombre positif = rotations normales (ra/rb),
// nombre negatif = rotations inverses (rra/rrb).
type move struct {
	src  int // rotations de la pile source (pour amener l'element au sommet)
	dst  int // rotations de la pile destination (pour preparer sa place)
	cost int // nombre total d'instructions de rotation
}

// cheapestMove parcourt tous les elements de src et renvoie le deplacement
// qui coute le moins d'instructions. target(dst, v) donne l'indice de
// l'element de dst qui doit etre au sommet pour que v soit bien place.
func cheapestMove(src, dst []int, target func([]int, int) int) move {
	best := move{cost: -1}
	for i, v := range src {
		m := bestRotation(i, len(src), target(dst, v), len(dst))
		if best.cost < 0 || m.cost < best.cost {
			best = m
		}
	}
	return best
}

// bestRotation : on doit amener l'indice i de la source (taille ls) et
// l'indice j de la destination (taille ld) au sommet de leur pile.
// Il y a 4 facons de faire, on garde la moins chere :
//
//   - les deux vers le haut     : max(i, j)   grace a "rr"
//   - les deux vers le bas      : max(ls-i, ld-j) grace a "rrr"
//   - source haut, dest bas     : i + (ld-j)
//   - source bas, dest haut     : (ls-i) + j
func bestRotation(i, ls, j, ld int) move {
	srcUp, srcDown := i, i-ls // srcDown est negatif
	dstUp, dstDown := j, j-ld // dstDown est negatif

	candidates := []move{
		{srcUp, dstUp, max(srcUp, dstUp)},
		{srcDown, dstDown, max(-srcDown, -dstDown)},
		{srcUp, dstDown, srcUp - dstDown},
		{srcDown, dstUp, -srcDown + dstUp},
	}
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.cost < best.cost {
			best = c
		}
	}
	return best
}

// applyMove effectue les rotations : d'abord les rotations communes (rr ou
// rrr) tant que les deux piles tournent dans le meme sens, puis le reste
// pile par pile. rotA et rotB suivent la convention de move.
func applyMove(r *recorder, rotA, rotB int) {
	for rotA > 0 && rotB > 0 {
		r.do("rr", 1)
		rotA, rotB = rotA-1, rotB-1
	}
	for rotA < 0 && rotB < 0 {
		r.do("rrr", 1)
		rotA, rotB = rotA+1, rotB+1
	}
	rotateTimes(r, "ra", "rra", rotA)
	rotateTimes(r, "rb", "rrb", rotB)
}

// rotateTimes fait n rotations : n > 0 avec fwd, n < 0 avec rev.
func rotateTimes(r *recorder, fwd, rev string, n int) {
	if n > 0 {
		r.do(fwd, n)
	} else {
		r.do(rev, -n)
	}
}

// shortestRotation renvoie le nombre de rotations (positif ou negatif) pour
// amener l'indice i au sommet d'une pile de taille n par le chemin le plus court.
func shortestRotation(i, n int) int {
	if i <= n/2 {
		return i
	}
	return i - n
}

// targetInDescending : b est triee en ordre decroissant. Pour inserer v, il
// faut qu'au sommet de b se trouve le plus grand element plus petit que v
// (v se posera juste au-dessus). Si v est le plus petit de tous, il doit se
// poser au-dessus du maximum de b (la "jonction" de l'ordre circulaire).
func targetInDescending(b []int, v int) int {
	bestIdx := -1
	for i, x := range b {
		if x < v && (bestIdx < 0 || x > b[bestIdx]) {
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return indexOfMax(b)
	}
	return bestIdx
}

// targetInAscending : symetrique pour a (ordre croissant). Au sommet de a
// doit se trouver le plus petit element plus grand que v. Si v est le plus
// grand, il se pose au-dessus du minimum de a.
func targetInAscending(a []int, v int) int {
	bestIdx := -1
	for i, x := range a {
		if x > v && (bestIdx < 0 || x < a[bestIdx]) {
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return indexOfMin(a)
	}
	return bestIdx
}

// sortThree trie une pile a de 3 elements en 2 instructions maximum :
//   - si le maximum est au sommet, "ra" l'envoie en bas ;
//   - s'il est au milieu, "rra" l'envoie en bas ;
//   - le maximum etant en bas, il reste au plus un "sa" a faire.
func sortThree(r *recorder) {
	switch indexOfMax(r.s.A) {
	case 0:
		r.do("ra", 1)
	case 1:
		r.do("rra", 1)
	}
	if r.s.A[0] > r.s.A[1] {
		r.do("sa", 1)
	}
}

// indexOfMin renvoie l'indice du plus petit element.
func indexOfMin(s []int) int {
	idx := 0
	for i, v := range s {
		if v < s[idx] {
			idx = i
		}
	}
	return idx
}

// indexOfMax renvoie l'indice du plus grand element.
func indexOfMax(s []int) int {
	idx := 0
	for i, v := range s {
		if v > s[idx] {
			idx = i
		}
	}
	return idx
}
