package connections

import (
	"fmt"
	"log"
	"os"
	"tradeapi/database/migrations"
	"tradeapi/database/seeds"
	"tradeapi/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

// GetDatabaseConnection cria e retorna uma conexão com o banco de dados usando o GORM
func GetDatabaseConnection(app *fiber.App, ctx *cli.Context, cmdType string) error {
	// Obter configurações do ambiente
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// String de conexão
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Cria a conexão
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	// Testar a conexão (sem Ping)
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter DB do GORM: %v", err)
	}

	// Testa a conexão com o banco
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("erro ao verificar conexão com o banco de dados: %v", err)
	}

	r := routes.NewRoute(db, app)
	r.SetupAPIRoutes()
	executeCommands(db, cmdType, ctx)
	
	return nil
}


func executeCommands(db *gorm.DB, cmdType string, ctx *cli.Context) {
    switch cmdType {
    case "migrate":
        m := migrations.NewMigration(db)
        m.RunMigrations()
		break
    case "seed":
        s := seeds.NewSeeder(db)
        model := ctx.String("model")
        if model == "" {
            log.Println("Rodando todas as seeders...")
            s.RunAllSeeds()
        } else {
            s.RunModelSeed(model)
        }
		break
	}
}