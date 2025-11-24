package main

import (
	"github.com/William-Le-Gavrian/go-projet-final/cmd"
	_ "github.com/William-Le-Gavrian/go-projet-final/cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "github.com/William-Le-Gavrian/go-projet-final/cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
)

func main() {
	cmd.Execute()
}
