// Commande checker : lit des instructions sur l'entree standard, les execute
// sur la pile construite a partir des arguments, puis affiche OK si a est
// triee et b vide, KO sinon.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"push-swap/internal/parser"
	"push-swap/internal/stack"
)

func main() {
	values, err := parser.ParseArgs(os.Args[1:])
	if err != nil {
		exitWithError()
	}

	// Sans argument : on s'arrete sans rien afficher (et sans lire stdin).
	if len(values) == 0 {
		return
	}

	s := stack.New(values)

	// On lit l'entree standard ligne par ligne : une instruction par ligne.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Sous Windows les fins de ligne peuvent etre "\r\n" : on retire le
		// "\r" pour ne pas rejeter une instruction valide.
		line := strings.TrimSuffix(scanner.Text(), "\r")

		// Une ligne vide est ignoree : `echo -e "sa\n"` ajoute une ligne vide
		// a la fin, ce n'est pas une instruction.
		if line == "" {
			continue
		}

		// Une instruction inconnue ou mal formatee (" sa", "SA", "sa ")
		// est refusee par Apply -> Error.
		if err := s.Apply(line); err != nil {
			exitWithError()
		}
	}
	if err := scanner.Err(); err != nil {
		exitWithError()
	}

	if s.IsSolved() {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

// exitWithError affiche "Error\n" sur stderr et quitte avec un code d'erreur.
func exitWithError() {
	fmt.Fprintln(os.Stderr, "Error")
	os.Exit(1)
}
