// Package stack represente les deux piles a et b et implemente les 11
// instructions du sujet.
//
// Representation choisie : une pile est un slice d'entiers ou l'indice 0 est
// le SOMMET de la pile. Ainsi a[0] est l'element que "sa", "pb" ou "ra"
// manipulent en premier, et l'affichage du sujet (de haut en bas) correspond
// directement a l'ordre du slice.
//
// Toutes les operations sont centralisees ici : push-swap (pour simuler son
// tri) et checker (pour executer les instructions lues) utilisent exactement
// le meme code, ce qui garantit qu'ils sont coherents.
package stack

import "fmt"

// Ops liste toutes les instructions valides. Le solveur en BFS l'utilise pour
// essayer chaque instruction possible.
var Ops = []string{"sa", "sb", "ss", "pa", "pb", "ra", "rb", "rr", "rra", "rrb", "rrr"}

// Stacks regroupe les deux piles du probleme.
type Stacks struct {
	A []int // pile a : contient les nombres au depart
	B []int // pile b : vide au depart, sert de pile auxiliaire
}

// New cree les piles a partir des valeurs recues. On copie les valeurs pour
// ne jamais modifier le slice de l'appelant par effet de bord.
func New(values []int) *Stacks {
	a := make([]int, len(values))
	copy(a, values)
	return &Stacks{A: a, B: []int{}}
}

// Clone renvoie une copie independante des deux piles (utile au BFS qui doit
// explorer plusieurs etats a partir d'un meme etat).
func (s *Stacks) Clone() *Stacks {
	c := &Stacks{A: make([]int, len(s.A)), B: make([]int, len(s.B))}
	copy(c.A, s.A)
	copy(c.B, s.B)
	return c
}

// Apply execute une instruction. Une instruction inconnue renvoie une erreur
// (le checker affichera alors "Error").
//
// Les instructions sur une pile vide ou trop petite (ex : "sa" avec un seul
// element, "pa" avec b vide) ne font simplement rien : c'est le comportement
// attendu et le programme ne plante jamais.
func (s *Stacks) Apply(op string) error {
	switch op {
	case "sa":
		swap(s.A)
	case "sb":
		swap(s.B)
	case "ss":
		swap(s.A)
		swap(s.B)
	case "pa":
		s.B, s.A = push(s.B, s.A)
	case "pb":
		s.A, s.B = push(s.A, s.B)
	case "ra":
		rotate(s.A)
	case "rb":
		rotate(s.B)
	case "rr":
		rotate(s.A)
		rotate(s.B)
	case "rra":
		reverseRotate(s.A)
	case "rrb":
		reverseRotate(s.B)
	case "rrr":
		reverseRotate(s.A)
		reverseRotate(s.B)
	default:
		return fmt.Errorf("unknown instruction %q", op)
	}
	return nil
}

// IsSolved indique si le but est atteint : a triee en ordre croissant
// (le plus petit au sommet) ET b vide.
func (s *Stacks) IsSolved() bool {
	return len(s.B) == 0 && IsSorted(s.A)
}

// IsSorted indique si les valeurs sont en ordre croissant.
// Une liste vide ou d'un seul element est consideree comme triee.
func IsSorted(values []int) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] {
			return false
		}
	}
	return true
}

// swap echange les deux premiers elements (sa / sb).
func swap(s []int) {
	if len(s) < 2 {
		return
	}
	s[0], s[1] = s[1], s[0]
}

// push retire le sommet de "from" et le place au sommet de "to" (pa / pb).
// Les slices changent de longueur, donc on renvoie les deux nouveaux slices.
func push(from, to []int) ([]int, []int) {
	if len(from) == 0 {
		return from, to
	}
	top := from[0]
	from = from[1:]

	// On agrandit "to" d'une case, on decale tout d'un cran vers le bas,
	// puis on place l'element au sommet (indice 0).
	to = append(to, 0)
	copy(to[1:], to)
	to[0] = top
	return from, to
}

// rotate : le premier element devient le dernier (ra / rb).
// Exemple : [1 2 3] -> [2 3 1]
func rotate(s []int) {
	if len(s) < 2 {
		return
	}
	top := s[0]
	copy(s, s[1:])
	s[len(s)-1] = top
}

// reverseRotate : le dernier element devient le premier (rra / rrb).
// Exemple : [1 2 3] -> [3 1 2]
func reverseRotate(s []int) {
	if len(s) < 2 {
		return
	}
	bottom := s[len(s)-1]
	copy(s[1:], s)
	s[0] = bottom
}
