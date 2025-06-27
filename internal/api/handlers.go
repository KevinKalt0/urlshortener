package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/KevinKalt0/urlshortener/internal/models"
	"github.com/KevinKalt0/urlshortener/internal/services"
	"github.com/KevinKalt0/urlshortener/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TODO Créer une variable ClickEventsChannel qui est un chan de type ClickEvent
var ClickEventsChannel chan services.ClickEvent

// SetupRoutes configure toutes les routes de l'API Gin et injecte les dépendances nécessaires
func SetupRoutes(router *gin.Engine, linkService *services.LinkService) {
	// Le channel est initialisé ici.
	if ClickEventsChannel == nil {
		// TODO Créer le channel ici (make), il doit être bufférisé
		// La taille du buffer doit être configurable via Viper (cfg.Analytics.BufferSize)
		cfg := config.GetConfig()
		ClickEventsChannel = make(chan services.ClickEvent, cfg.Analytics.BufferSize)
	}

	// TODO : Route de Health Check , /health
	router.GET("/health", HealthCheckHandler)

	// TODO : Routes de l'API
	// Doivent être au format /api/v1/
	// POST /links
	// GET /links/:shortCode/stats
	api := router.Group("/api/v1")
	{
		api.POST("/links", CreateShortLinkHandler(linkService))
		api.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))
	}

	// Route de Redirection (au niveau racine pour les short codes)
	router.GET("/:shortCode", RedirectHandler(linkService))
}

// HealthCheckHandler gère la route /health pour vérifier l'état du service.
func HealthCheckHandler(c *gin.Context) {
	// TODO  Retourner simplement du JSON avec un StatusOK, {"status": "ok"}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateLinkRequest représente le corps de la requête JSON pour la création d'un lien.
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

// CreateShortLinkHandler gère la création d'une URL courte.
func CreateShortLinkHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		// TODO : Tente de lier le JSON de la requête à la structure CreateLinkRequest.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format or missing long_url"})
			return
		}

		// TODO: Appeler le LinkService (CreateLink pour créer le nouveau lien.
		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short link"})
			return
		}

		cfg := config.GetConfig()

		// Retourne le code court et l'URL longue dans la réponse JSON.
		// TODO Choisir le bon code HTTP
		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": cfg.Server.BaseURL + "/" + link.ShortCode,
		})
	}
}

// RedirectHandler gère la redirection d'une URL courte vers l'URL longue et l'enregistrement asynchrone des clics.
func RedirectHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le shortCode de l'URL avec c.Param
		shortCode := c.Param("shortCode")

		// TODO 2: Récupérer l'URL longue associée au shortCode depuis le linkService (GetLinkByShortCode)
		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
				return
			}
			log.Printf("Error retrieving link for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// TODO 3: Créer un ClickEvent avec les informations pertinentes.
		clickEvent := services.ClickEvent{
			LinkID:    link.ID,
			Timestamp: time.Now(),
			UserAgent: c.Request.UserAgent(),
			IP:        c.ClientIP(),
		}

		// TODO 4: Envoyer le ClickEvent dans le ClickEventsChannel avec le Multiplexage.
		select {
		case ClickEventsChannel <- clickEvent:
			// OK envoyé
		default:
			log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		}

		// TODO 5: Effectuer la redirection HTTP 302 (StatusFound) vers l'URL longue.
		c.Redirect(http.StatusFound, link.LongURL)
	}
}

// GetLinkStatsHandler gère la récupération des statistiques pour un lien spécifique.
func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO Récupère le shortCode de l'URL avec c.Param
		shortCode := c.Param("shortCode")

		// TODO 6: Appeler le LinkService pour obtenir le lien et le nombre total de clics.
		link, totalClicks, err := linkService.GetStatsByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
			return
		}

		// Retourne les statistiques dans la réponse JSON.
		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
