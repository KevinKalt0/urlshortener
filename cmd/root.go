package cmd

import (
	"github.com/KevinKalt0/urlshortener/internal/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Cfg est la variable globale qui contiendra la configuration chargée.
// Elle sera accessible à toutes les commandes Cobra.
var Cfg *config.Config

// RootCmd représente la commande de base lorsque l'on appelle l'application sans sous-commande.
// C'est le point d'entrée principal pour Cobra.
var RootCmd = &cobra.Command{
	Use:   "urlshortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: `'url-shortener' est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.
Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

// Execute est le point d'entrée principal pour l'application Cobra.
// Il est appelé depuis 'main.go'.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Charger la configuration au démarrage de l'application
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {
		// En cas d'erreur de chargement de config, on peut soit :
		// 1. Arrêter l'application
		// 2. Utiliser des valeurs par défaut
		// 3. Logger l'erreur et continuer

		// Option 1 : Arrêter si la config ne peut pas être chargée
		log.Printf("Attention: Impossible de charger la configuration (%v), utilisation des valeurs par défaut", err)

		// Option 2 : Continuer avec des valeurs par défaut (recommandé pour le développement)
		// log.Fatal("Impossible de charger la configuration:", err)
	}

	// Les sous-commandes seront ajoutées ici une fois qu'elles seront créées
	// Décommentez ces lignes quand vous aurez créé les fichiers correspondants :
	// RootCmd.AddCommand(migrateCmd)
	// RootCmd.AddCommand(runServerCmd)
	// RootCmd.AddCommand(createCmd)
	// RootCmd.AddCommand(statsCmd)

	// Flags persistants (optionnel)
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (défaut: ./configs/config.yaml)")
}
