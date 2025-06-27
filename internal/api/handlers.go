package api

import (
    "net/http"
    "github.com/KevinKalt0/internal/models"
    "github.com/KevinKalt0/internal/repository"

<<<<<<< Updated upstream
	"github.com/KevinKalt0/urlshortener/internal/models"
	"github.com/KevinKalt0/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Pour gérer gorm.ErrRecordNotFound
=======
    "github.com/gin-gonic/gin"
>>>>>>> Stashed changes
)

type APIHandler struct {
    LinkRepo  repository.LinkRepository
    ClickRepo repository.ClickRepository
}

// Constructeur de handler
func NewAPIHandler(linkRepo repository.LinkRepository, clickRepo repository.ClickRepository) *APIHandler {
    return &APIHandler{
        LinkRepo:  linkRepo,
        ClickRepo: clickRepo,
    }
}

// GET /health
func (h *APIHandler) HealthCheckHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// POST /api/v1/links
func (h *APIHandler) CreateShortLinkHandler(c *gin.Context) {
    var req struct {
        LongURL string `json:"long_url" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Générer un shortcode de 6 caractères (logique à placer dans link_service.go)
    shortCode := generateShortCode() // Placeholder - à remplacer par appel réel au service
    link := &models.Link{
        ShortCode: shortCode,
        LongURL:   req.LongURL,
    }

    if err := h.LinkRepo.Create(link); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save link"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "code":      link.ShortCode,
        "short_url": "http://localhost:8080/" + link.ShortCode,
    })
}

// GET /:shortCode
func (h *APIHandler) RedirectHandler(c *gin.Context) {
    code := c.Param("shortCode")
    link, err := h.LinkRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
        return
    }

    // Ici on enverra un événement au worker pour enregistrer le clic (non bloquant)
    // Exemple : clickChannel <- ClickEvent{LinkID: link.ID, IP: c.ClientIP(), UA: ...}

    c.Redirect(http.StatusFound, link.LongURL)
}

// GET /api/v1/links/:shortCode/stats
func (h *APIHandler) GetLinkStatsHandler(c *gin.Context) {
    code := c.Param("shortCode")
    link, err := h.LinkRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
        return
    }

    count, err := h.ClickRepo.CountByLinkID(link.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stats"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "short_code": link.ShortCode,
        "long_url":   link.LongURL,
        "clicks":     count,
    })
}
