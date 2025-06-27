package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/KevinKalt0/urlshortener/cmd"
	"github.com/KevinKalt0/urlshortener/internal/repository"
	"github.com/KevinKalt0/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 1: Valider que le flag --url a été fourni.
		if _, err := url.ParseRequestURI(longURLFlag); err != nil {
			fmt.Println("L'URL fournie n'est pas valide.")
			os.Exit(1)
		}
		// TODO Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		
		// TODO : Initialiser la connexion à la base de données SQLite.

		db, err := gorm.Open(sqlite.Open("url_shortener.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Impossible de se connecter à la base : %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Impossible d'obtenir l'instance SQL : %v", err)
		}
		defer sqlDB.Close()

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)
		linkService := services.NewLinkService(linkRepo, clickRepo)
		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		// os.Exit(1) si erreur
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Fatalf("FATAL: Erreur lors de la création du lien court : %v", err)
		}
		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// TODO : Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
	cmd2.RootCmd.AddCommand(CreateCmd)
	// TODO :  Marquer le flag comme requis

	// TODO : Ajouter la commande à RootCmd

}
