package solver

import (
	"math/rand"
	"testing"

	"push-swap/internal/stack"
)

// check execute ops sur values et verifie que le resultat est trie
// (c'est exactement ce que fait le programme checker).
func check(t *testing.T, values []int, ops []string) {
	t.Helper()
	s := stack.New(values)
	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			t.Fatalf("instruction invalide %q", op)
		}
	}
	if !s.IsSolved() {
		t.Fatalf("%v n'est pas trie apres %d instructions", values, len(ops))
	}
}

func TestSubjectExamples(t *testing.T) {
	tests := []struct {
		values []int
		maxOps int
	}{
		{[]int{2, 1, 3, 6, 5, 8}, 7}, // le sujet donne 7 instructions
		{[]int{4, 67, 3, 87, 23}, 6}, // le sujet donne 6 instructions
	}
	for _, tt := range tests {
		ops := Solve(tt.values)
		check(t, tt.values, ops)
		if len(ops) > tt.maxOps {
			t.Errorf("%v : %d instructions, attendu au plus %d", tt.values, len(ops), tt.maxOps)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	for _, values := range [][]int{nil, {42}, {1, 2, 3, 4, 5, 6, 7, 8}} {
		if ops := Solve(values); len(ops) != 0 {
			t.Errorf("%v deja triee : %d instructions au lieu de 0", values, len(ops))
		}
	}
	check(t, []int{2, 1}, Solve([]int{2, 1}))
}

// TestAllPermutationsSmall teste toutes les permutations de 1 a 6 elements.
func TestAllPermutationsSmall(t *testing.T) {
	for n := 1; n <= 6; n++ {
		forEachPermutation(n, func(p []int) {
			check(t, p, Solve(p))
		})
	}
}

// TestRandomLarge teste des listes aleatoires de grande taille et affiche
// le nombre moyen d'instructions (visible avec go test -v).
func TestRandomLarge(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for _, n := range []int{7, 10, 50, 100, 500} {
		total, worst := 0, 0
		const runs = 20
		for i := 0; i < runs; i++ {
			values := rng.Perm(n)
			ops := Solve(values)
			check(t, values, ops)
			total += len(ops)
			worst = max(worst, len(ops))
		}
		t.Logf("n=%d : moyenne %d instructions, pire cas %d", n, total/runs, worst)
	}
}

// forEachPermutation appelle f pour chaque permutation de 0..n-1.
func forEachPermutation(n int, f func([]int)) {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	var permute func(k int)
	permute = func(k int) {
		if k == n {
			f(append([]int(nil), p...))
			return
		}
		for i := k; i < n; i++ {
			p[k], p[i] = p[i], p[k]
			permute(k + 1)
			p[k], p[i] = p[i], p[k]
		}
	}
	permute(0)
}
