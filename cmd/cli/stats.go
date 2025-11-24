package cli

import (
	"fmt"
	"log"
	"os"

	cmd2 "github.com/William-Le-Gavrian/go-projet-final/cmd"
	"github.com/William-Le-Gavrian/go-projet-final/internal/models"
	"github.com/William-Le-Gavrian/go-projet-final/internal/repository"
	"github.com/William-Le-Gavrian/go-projet-final/internal/services"
	"github.com/spf13/cobra"

	"github.com/glebarez/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : variable shortCodeFlag qui stockera la valeur du flag --code

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		codeFlag, _ := cmd.Flags().GetString("code")

		// TODO : Valider que le flag --code a été fourni.
		// os.Exit(1) si erreur
		if codeFlag == "" {
			log.Println("Erreur: paramètre --code requis")
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("Erreur en chargeant la configuration")
		}

		// TODO 3: Initialiser la connexion à la BDD.
		// log.Fatalf si erreur
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

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande grâce à defer
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// TODO 5: Appeler GetLinkStats pour récupérer le lien et ses statistiques.
		// Attention, la fonction retourne 3 valeurs
		// Pour l'erreur, utilisez gorm.ErrRecordNotFound
		// Si erreur, os.Exit(1)

		link, totalClicks, err := linkService.GetLinkStats(codeFlag)
		if err != nil {
			log.Printf("Erreur lors de la récupération des statistiques: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.Shortcode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// TODO : Définir le flag --code pour la commande create.
	StatsCmd.Flags().StringP("url", "u", "", "Code pour récuperer les statistiques")

	// TODO :  Marquer le flag comme requis
	StatsCmd.MarkFlagRequired("code")

	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(StatsCmd)

}
