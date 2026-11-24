package parser

import (
	"reflect"
	"testing"
)

func TestParseArgsValid(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []int
	}{
		{"aucun argument", nil, nil},
		{"un seul argument avec espaces", []string{"2 1 3"}, []int{2, 1, 3}},
		{"plusieurs arguments", []string{"2", "1", "3"}, []int{2, 1, 3}},
		{"espaces multiples", []string{"  4   -2 "}, []int{4, -2}},
		{"argument vide", []string{""}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if err != nil {
				t.Fatalf("erreur inattendue : %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgs(%q) = %v, attendu %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestParseArgsInvalid(t *testing.T) {
	tests := [][]string{
		{"0 one 2 3"},            // pas un entier
		{"1 2 1"},                // doublon
		{"1", "1"},               // doublon sur deux arguments
		{"1.5"},                  // decimal
		{"12abc"},                // caracteres en trop
		{"99999999999999999999"}, // depasse la taille d'un int
	}
	for _, args := range tests {
		if _, err := ParseArgs(args); err == nil {
			t.Errorf("ParseArgs(%q) aurait du renvoyer une erreur", args)
		}
	}
}
