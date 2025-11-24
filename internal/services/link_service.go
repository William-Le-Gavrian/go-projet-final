package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/William-Le-Gavrian/go-projet-final/internal/models"
	"github.com/William-Le-Gavrian/go-projet-final/internal/repository" // Importe le package repository
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
	var shortCode = make([]byte, length)
	lenCharset := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		randInt, err := rand.Int(rand.Reader, lenCharset)
		if err != nil {
			return "", err
		}
		shortCode = append(shortCode, charset[randInt.Int64()])
	}

	return string(shortCode), nil
}

// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	// TODO 1: Implémenter la logique de retry pour générer un code court unique.
	// Essayez de générer un code, vérifiez s'il existe déjà en base, et retentez si une collision est trouvée.
	// Limitez le nombre de tentatives pour éviter une boucle infinie.

	// TODO Créer une variable shortcode pour stocker le shortcode créé

	// TODO Définir un nombre maximum (5) de tentative pour trouver un code unique  (maxRetries)
	maxRetries := 5

	var shortCode string

	for i := 0; i < maxRetries; i++ {
		// TODO : Génère un code de 6 caractères (GenerateShortCode)
		code, err := s.GenerateShortCode(6)

		// TODO : Vérifie si le code généré existe déjà en base de données (GetLinkbyShortCode)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code
				break
			}
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}

		// Si aucune erreur (le code a été trouvé), cela signifie une collision.
		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)
		// La boucle continuera pour générer un nouveau code.
	}

	// TODO : Si après toutes les tentatives, aucun code unique n'a été trouvé... Errors.New
	if shortCode == "" {
		return nil, errors.New("failed to generate a unique short code after maximum retries")
	}

	// TODO Crée une nouvelle instance du modèle Link.
	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreateAt:  time.Now(),
	}

	// TODO Persiste le nouveau lien dans la base de données via le repository (CreateLink)
	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	// TODO Retourne le lien créé
	return link, nil

}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)

	if err != nil {
		return nil, fmt.Errorf("failed to get link by short code %s: %w", shortCode, err)
	}

	return link, nil
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// TODO : Récupérer le lien par son shortCode
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get link by short code %s: %w", shortCode, err)
	}

	// TODO 4: Compter le nombre de clics pour ce LinkID
	count, err := s.linkRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks for link ID %d: %w", link.ID, err)
	}

	// TODO : on retourne les 3 valeurs
	return link, count, nil
}
