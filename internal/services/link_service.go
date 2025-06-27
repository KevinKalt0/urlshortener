package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/KevinKalt0/urlshortener/internal/models"
	"github.com/KevinKalt0/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// TODO Créer la struct
// LinkService est une structure qui g fournit des méthodes pour la logique métier des liens.
// Elle détient linkRepo qui est une référence vers une interface LinkRepository.
// IMPORTANT : Le champ doit être du type de l'interface (non-pointeur).
type LinkService struct {
	linkRepo repository.LinkRepository
}

// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// TODO Créer la méthode GenerateShortCode
// GenerateShortCode est une méthode rattachée à LinkService
// Elle génère un code court aléatoire d'une longueur spécifiée. Elle prend une longueur en paramètre et retourne une string et une erreur
// Il utilise le package 'crypto/rand' pour éviter la prévisibilité.
// Je vous laisse chercher un peu :) C'est faisable en une petite dizaine de ligne
func (s *LinkService) GenerateShortCode(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	// TODO 1: Implémenter la logique de retry pour générer un code court unique.
	// Essayez de générer un code, vérifiez s'il existe déjà en base, et retentez si une collision est trouvée.
	// Limitez le nombre de tentatives pour éviter une boucle infinie.

	// TODO Créer une variable shortcode pour stocker le shortcode créé
	var shortCode string

	// TODO Définir un nombre maximum (5) de tentative pour trouver un code unique  (maxRetries)
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		// TODO : Génère un code de 6 caractères (GenerateShortCode)
		code, err := s.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// TODO : Vérifie si le code généré existe déjà en base de données (GetLinkbyShortCode)
		// On ignore la première valeur
		_, err = s.linkRepo.GetLinkByShortCode(code)

		if err != nil {
			// Si l'erreur est 'record not found' de GORM, cela signifie que le code est unique.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code // Le code est unique, on peut l'utiliser
				break            // Sort de la boucle de retry
			}
			// Si c'est une autre erreur de base de données, retourne l'erreur.
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}

		// Si aucune erreur (le code a été trouvé), cela signifie une collision.
		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)
		// La boucle continuera pour générer un nouveau code.
	}

	// TODO : Si après toutes les tentatives, aucun code unique n'a été trouvé... Errors.New
	if shortCode == "" {
		return nil, errors.New("failed to generate unique short code after maximum retries")
	}

	// TODO Crée une nouvelle instance du modèle Link.
	link := &models.Link{
		ShortCode:  shortCode,
		LongURL:    longURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ClickCount: 0,
	}

	// TODO Persiste le nouveau lien dans la base de données via le repository (CreateLink)
	err := s.linkRepo.CreateLink(link)
	if err != nil {
		return nil, fmt.Errorf("failed to create link in database: %w", err)
	}

	// TODO Retourne le lien créé
	return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	// TODO : Récupérer un lien par son code court en utilisant s.linkRepo.GetLinkByShortCode.
	// Retourner le lien trouvé ou une erreur si non trouvé/problème DB.
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by short code: %w", err)
	}
	return link, nil
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// TODO : Récupérer le lien par son shortCode
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get link by short code: %w", err)
	}

	// TODO 4: Compter le nombre de clics pour ce LinkID
	// Note: Ici on utilise le ClickCount du modèle Link qui devrait être maintenu à jour
	// Alternativement, si vous avez un ClickRepository, vous pourriez faire :
	// clickCount, err := s.clickRepo.CountClicksByLinkID(link.ID)
	clickCount := int(link.ClickCount)

	// TODO : on retourne les 3 valeurs
	return link, clickCount, nil
}
