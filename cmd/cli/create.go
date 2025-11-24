package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/William-Le-Gavrian/go-projet-final/cmd"
	"github.com/William-Le-Gavrian/go-projet-final/internal/models"
	"github.com/William-Le-Gavrian/go-projet-final/internal/repository"
	"github.com/William-Le-Gavrian/go-projet-final/internal/services"
	"github.com/glebarez/sqlite" // Driver SQLite pour GORM
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		longURL, _ := cmd.Flags().GetString("url")

		// TODO 1: Valider que le flag --url a été fourni.
		if longURL == "" {
			log.Println("Erreur: paramètre --url requis")
			os.Exit(1)
		}

		// TODO Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)
		_, err := url.Parse(longURL)
		if err != nil {
			log.Printf("Erreur: URL invalide: %v", err)
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("Erreur en chargeant la configuration")
		}

		// TODO : Initialiser la connexion à la base de données SQLite.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		err = db.AutoMigrate(&models.Link{}, &models.Click{})
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		// os.Exit(1) si erreur
		link, err := linkService.CreateLink(longURL)
		if err != nil {
			log.Printf("Erreur lors de la création du lien: %v", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.Shortcode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.Shortcode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// TODO : Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringP("url", "u", "", "URL longue à raccourcir")

	// TODO :  Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")

	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)

}
