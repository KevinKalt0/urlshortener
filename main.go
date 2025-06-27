/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/KevinKalt0/urlshortener/cmd"
	_ "github.com/KevinKalt0/urlshortener/cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "github.com/KevinKalt0/urlshortener/cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
)

func main() {
	// TODO Exécute la commande racine de Cobra.
	cmd.Execute()

}
