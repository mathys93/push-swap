// Commande push-swap : affiche sur la sortie standard la liste d'instructions
// qui trie les entiers recus en argument.
//
// Ce fichier ne fait que l'orchestration et l'affichage :
// parsing -> calcul du tri -> affichage. Toute la logique est dans internal/.
package main

import (
	"bufio"
	"fmt"
	"os"

	"push-swap/internal/parser"
	"push-swap/internal/solver"
)

func main() {
	// os.Args[0] est le nom du programme, les arguments commencent a 1.
	values, err := parser.ParseArgs(os.Args[1:])
	if err != nil {
		// Le sujet impose "Error\n" sur la sortie d'erreur (stderr).
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}

	// Sans argument, values est vide : Solve ne renvoie rien et le programme
	// n'affiche rien, comme demande.
	//
	// bufio.Writer regroupe les ecritures : pour 500 nombres il y a des
	// milliers de lignes, c'est bien plus rapide qu'un appel systeme par ligne.
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for _, op := range solver.Solve(values) {
		fmt.Fprintln(out, op) // une instruction par ligne, separees par \n
	}
}
