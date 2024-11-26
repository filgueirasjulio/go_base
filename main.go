package main

import (
	"fmt"
	"os"
	"log"
	"os/signal"
	"syscall"
	"time"
	"gorm.io/gorm"

	"tradeapi/connections"
	"tradeapi/routes"
	seeds "tradeapi/database/seeds"
	migrations "tradeapi/database/migrations"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func init() {
	// Carrega as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env:", err)
	}
}

func main() {
	cmd := &cli.App{
		Name:        "Tradeapi",
		Usage:       "Versão de Desenvolvimento",
		Version:     "1.0.0",
		UsageText:   "Rode main [global options] command [command options] [arguments...]",
		Description: "Api de trade",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Executando em modo padrão",
				Action: func(ctx *cli.Context) error {
					// Conectar ao banco e rodar o serviço
					db := Teste()

					Run(db)
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Executando teste de conexão",
				Action: func(ctx *cli.Context) error {
					err := Teste()			
					if err != nil {
						fmt.Println("Falha no teste de conexão:", err)
					}

					fmt.Println("Teste de conexão bem-sucedido!")
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "Executar as migrations",
				Action: func(ctx *cli.Context) error {
					db := Teste()

					migrations.RunMigrations(db)
					return nil
				},
			},
			{
				Name:  "seed",
				Usage: "Executa as seeders",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "model",
						Aliases:  []string{"m"},
						Usage:    "Modelo para seed (ex: User)",
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					// Estabelecendo a conexão com o banco
					db, err := connections.GetDatabaseConnection()
					if err != nil {
						log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
					}

					// Verifica qual modelo foi passado na flag
					model := c.String("model")
					if model == "" {
						log.Println("Rodando todas as seeders...")
						seeds.RunAllSeeds(db)
					} else {
						seeds.RunModelSeed(db, model)
					}

					return nil
				},
			},
		},
		Compiled:  time.Now(),
		Authors:   []*cli.Author{},
		Copyright: fmt.Sprintf("© %d Tradeapi - todos direitos reservados", time.Now().Year()),
	}

	cmd.Run(os.Args)
}

func Teste() *gorm.DB {
	db, err := connections.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return db
}

func Run(db *gorm.DB) {
	// Crie uma instância do Fiber
	app := fiber.New()

	// Defina rotas
	routes.SetupAPIRoutes(app, db)

	// Defina porta
	port := os.Getenv("APP_PORT")

	go func() {
		fmt.Printf("Servidor rodando na porta %s...\n", port)
		if err := app.Listen(":"+port); err != nil {
			fmt.Println("Erro ao rodar o servidor:", err)
		}
	}()

	// Aguarde sinal de interrupção
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	if err := app.Shutdown(); err != nil {
        fmt.Println("Erro ao encerrar servidor:", err)
    }

	fmt.Println("Servidor encerrado!")
}
