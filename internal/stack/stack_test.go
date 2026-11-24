package stack

import (
	"reflect"
	"testing"
)

// TestSubjectExample rejoue l'exemple du sujet etape par etape.
func TestSubjectExample(t *testing.T) {
	s := New([]int{2, 1, 3, 6, 8, 5})
	steps := []struct {
		ops  []string
		a, b []int
	}{
		{[]string{"sa"}, []int{1, 2, 3, 6, 8, 5}, []int{}},
		{[]string{"pb", "pb", "pb"}, []int{6, 8, 5}, []int{3, 2, 1}},
		{[]string{"rb"}, []int{6, 8, 5}, []int{2, 1, 3}},
		{[]string{"rrr"}, []int{5, 6, 8}, []int{3, 2, 1}},
		{[]string{"pa", "pa", "pa"}, []int{1, 2, 3, 5, 6, 8}, []int{}},
	}
	for _, step := range steps {
		for _, op := range step.ops {
			if err := s.Apply(op); err != nil {
				t.Fatal(err)
			}
		}
		if !reflect.DeepEqual(s.A, step.a) || !reflect.DeepEqual(s.B, step.b) {
			t.Fatalf("apres %v : a=%v b=%v, attendu a=%v b=%v", step.ops, s.A, s.B, step.a, step.b)
		}
	}
	if !s.IsSolved() {
		t.Error("la pile devrait etre triee")
	}
}

// TestEachOperation verifie chaque instruction separement.
func TestEachOperation(t *testing.T) {
	tests := []struct {
		op    string
		wantA []int
		wantB []int
	}{
		{"sa", []int{2, 1, 3}, []int{5, 4, 6}},
		{"sb", []int{1, 2, 3}, []int{4, 5, 6}},
		{"ss", []int{2, 1, 3}, []int{4, 5, 6}},
		{"pa", []int{5, 1, 2, 3}, []int{4, 6}},
		{"pb", []int{2, 3}, []int{1, 5, 4, 6}},
		{"ra", []int{2, 3, 1}, []int{5, 4, 6}},
		{"rb", []int{1, 2, 3}, []int{4, 6, 5}},
		{"rr", []int{2, 3, 1}, []int{4, 6, 5}},
		{"rra", []int{3, 1, 2}, []int{5, 4, 6}},
		{"rrb", []int{1, 2, 3}, []int{6, 5, 4}},
		{"rrr", []int{3, 1, 2}, []int{6, 5, 4}},
	}
	for _, tt := range tests {
		s := &Stacks{A: []int{1, 2, 3}, B: []int{5, 4, 6}}
		if err := s.Apply(tt.op); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(s.A, tt.wantA) || !reflect.DeepEqual(s.B, tt.wantB) {
			t.Errorf("%s : a=%v b=%v, attendu a=%v b=%v", tt.op, s.A, s.B, tt.wantA, tt.wantB)
		}
	}
}

// TestEmptyStacks : les instructions sur des piles vides ne plantent pas.
func TestEmptyStacks(t *testing.T) {
	for _, op := range Ops {
		s := New(nil)
		if err := s.Apply(op); err != nil {
			t.Errorf("%s sur piles vides : %v", op, err)
		}
		if len(s.A) != 0 || len(s.B) != 0 {
			t.Errorf("%s ne devrait rien changer sur des piles vides", op)
		}
	}
}

func TestUnknownInstruction(t *testing.T) {
	for _, op := range []string{"", "foo", "SA", " sa", "sa ", "rrrr"} {
		if err := New([]int{1, 2}).Apply(op); err == nil {
			t.Errorf("%q aurait du etre refusee", op)
		}
	}
}

func TestIsSolved(t *testing.T) {
	if !(&Stacks{A: []int{1, 2, 3}}).IsSolved() {
		t.Error("a triee et b vide : devrait etre resolu")
	}
	if (&Stacks{A: []int{1, 3}, B: []int{2}}).IsSolved() {
		t.Error("b non vide : ne devrait pas etre resolu")
	}
	if (&Stacks{A: []int{2, 1}}).IsSolved() {
		t.Error("a non triee : ne devrait pas etre resolu")
	}
}
