package services

import "time"

import (
    "fmt"
    "github.com/KevinKalt0/urlshortener/internal/models"
    "github.com/KevinKalt0/urlshortener/internal/repository" // Importe le package repository
)

// TODO : créer la struct
// ClickService est une structure qui fournit des méthodes pour la logique métier des clics.
// Elle est juste composer de clickRepo qui est de type ClickRepository
type ClickService struct {
    clickRepo repository.ClickRepository
}
type ClickEvent struct {
	LinkID    uint
	Timestamp time.Time
	UserAgent string
	IP        string
}
// NewClickService crée et retourne une nouvelle instance de ClickService.
// C'est la fonction recommandée pour obtenir un service, assurant que toutes ses dépendances sont injectées.
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
    return &ClickService{
        clickRepo: clickRepo,
    }
}

// RecordClick enregistre un nouvel événement de clic dans la base de données.
// Cette méthode est appelée par le worker asynchrone.
func (s *ClickService) RecordClick(click *models.Click) error {
    // TODO 1: Appeler le ClickRepository (CreateClick) pour créer l'enregistrement de clic.
    // Gérer toute erreur provenant du repository.
    err := s.clickRepo.CreateClick(click)
    if err != nil {
        return fmt.Errorf("failed to record click: %w", err)
    }
    return nil
}

// GetClicksCountByLinkID récupère le nombre total de clics pour un LinkID donné.
// Cette méthode pourrait être utilisée par le LinkService pour les statistiques, ou directement par l'API stats.
func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
    // TODO 2: Appeler le ClickRepository (CountclicksByLinkID) pour compter les clics par LinkID.
    count, err := s.clickRepo.CountClicksByLinkID(linkID)
    if err != nil {
        return 0, fmt.Errorf("failed to get clicks count for link ID %d: %w", linkID, err)
    }
    return count, nil
}