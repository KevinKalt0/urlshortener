package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/KevinKalt0/urlshortener/cmd"
	"github.com/KevinKalt0/urlshortener/internal/api"
	"github.com/KevinKalt0/urlshortener/internal/models"
	"github.com/KevinKalt0/urlshortener/internal/monitor"
	"github.com/KevinKalt0/urlshortener/internal/repository"
	"github.com/KevinKalt0/urlshortener/internal/services"
	"github.com/KevinKalt0/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// RunServerCmd représente la commande 'run-server' de Cobra.
// C'est le point d'entrée pour lancer le serveur de l'application.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Charger la configuration chargée globalement via cmd.cfg
		// Ne pas oublier la gestion d'erreur (si nil ?), si erreur, faire un log.Fataf
		serverAddr := ":8080" // valeur hardcodée tant que cfg n’est pas utilisé

		// TODO : Initialiser la connexion à la base de données SQLite avec GORM.
		// Utilisez le nom de la base de données depuis la configuration (cfg.Database.Name).
		db, err := gorm.Open(sqlite.Open("url_shortener.db"), &gorm.Config{})
		if err != nil {
		log.Fatalf("Erreur de connexion à la base de données : %v", err)
		}

		// TODO : Initialiser les repositories.
		// Créez des instances de GormLinkRepository et GormClickRepository.

		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)

		// Laissez le log
		log.Println("Repositories initialisés.")

		// TODO : Initialiser les services métiers.
		// Créez des instances de LinkService et ClickService, en leur passant les repositories nécessaires.
		linkService := services.NewLinkService(linkRepo, clickRepo)
		// clickService := services.NewClickService(clickRepo)

		// Laissez le log
		log.Println("Services métiers initialisés.")

		// TODO : Initialiser le channel ClickEventsChannel (api/handlers) des événements de clic et lancer les workers (StartClickWorkers).
		// Le channel est bufferisé avec la taille configurée.
		// Passez le channel et le clickRepo aux workers.
		clickEvents := make(chan models.ClickEvent, 100) // buffer de 100

		go workers.StartClickWorkers(1, clickEvents, clickRepo)



		// TODO : Remplacer les XXX par les bonnes variables
		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
100, 1)

		// TODO : Initialiser et lancer le moniteur d'URLs.
		// Utilisez l'intervalle configuré (cfg.Monitor.IntervalMinutes).
		// Lancez le moniteur dans sa propre goroutine.
		monitorInterval := 5 * time.Minute // valeur par défaut en attendant Viper
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
		go urlMonitor.Start()
		log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

		// TODO : Configurer le routeur Gin et les handlers API.
		// Passez les services nécessaires aux fonctions de configuration des routes.
		// handler := api.SetupRoutes(linkService)
		// router := gin.Default()
		// router.GET("/health", handler.HealthCheckHandler)
		// router.POST("/api/v1/links", handler.CreateShortLinkHandler)
		// router.GET("/:shortCode", handler.RedirectHandler)
		// router.GET("/api/v1/links/:shortCode/stats", handler.GetLinkStatsHandler)
		router := gin.Default()

		clickEvents = make(chan services.ClickEvent, 100)

		api.SetupRoutes(router, linkService)


		// Pas toucher au log
		log.Println("Routes API configurées.")
		
		cfg := cmd2.Cfg

		// Créer le serveur HTTP Gin
		serverAddr = fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// TODO : Démarrer le serveur Gin dans une goroutine anonyme pour ne pas bloquer.
		// Pensez à logger des ptites informations...
		go func() {
			log.Printf("Serveur lancé sur http://localhost%s", serverAddr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erreur serveur : %v", err)
			}
			}()

		// Gére l'arrêt propre du serveur (graceful shutdown).
		// Créez un channel pour les signaux OS (SIGINT, SIGTERM).
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Attendre Ctrl+C ou signal d'arrêt

		// Bloquer jusqu'à ce qu'un signal d'arrêt soit reçu.
		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		// Arrêt propre du serveur HTTP avec un timeout.
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(RunServerCmd)
}

