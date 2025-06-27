/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// cfg est la variable globale qui contiendra la configuration chargée.
// Elle sera accessible à toutes les commandes Cobra.
var Cfg *config.Config

// TODO : Créer la RootCmd avec Cobra

// rootCmd représente la commande de base lorsque l'on appelle l'application sans sous-commande.
// C'est le point d'entrée principal pour Cobra.
var RootCmd = &cobra.Command{
	Use:   "URLshortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: "'url-shortener' est une application complète pour gérer des URLs courtes." +
		"\nElle inclut un serveur API pour le raccourcissement et la redirection,\n" +
		"ainsi qu'une interface en ligne de commande pour l'administration.\n" +
		"Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.",
}

// Execute est le point d'entrée principal pour l'application Cobra.
// Il est appelé depuis 'main.go'.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.urlshortener.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(runServerCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(statsCmd)
}


