package main

import (
	"fmt"
	"os"
	"log"
	"os/signal"
	"syscall"
	"strconv"
	"time"
	"tradeapi/connections"
	queue "tradeapi/app"
	
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)


var app *fiber.App

func init() {
	app = fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// Carrega as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env:", err)
	}

	//iniciar fila
	capacity, _ := strconv.Atoi(os.Getenv("QUEUE_CAPACITY"))

	queue.InitQueue(capacity)
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
					app, err := Teste(app, ctx, "run")
					if err != nil {
						log.Fatal("Erro ao conectar ao banco de dados: %v", err)
					}

					Run(app)
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Executando teste de conexão",
				Action: func(ctx *cli.Context) error {
					_, err := Teste(app, ctx, "test")			
					if err != nil {
						log.Fatal("Falha no teste de conexão:", err)
					}

					log.Println("Conexão com o banco de dados foi bem-sucedida!")
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "Executar as migrations",
				Action: func(ctx *cli.Context) error {
					_, err := Teste(app, ctx, "migrate")
					if err != nil {
						log.Fatal("Erro ao conectar ao banco de dados: %v", err)
					}

					log.Println("Conexão com o banco de dados foi bem-sucedida!")
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
				Action: func(ctx *cli.Context) error {
					// Estabelecendo a conexão com o banco
					_, err := Teste(app, ctx, "seed")
					if err != nil {
						log.Fatal("Erro ao conectar ao banco de dados: %v", err)
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

func Teste(app *fiber.App, c *cli.Context, cmdType string) (*fiber.App,  error) {
	err := connections.GetDatabaseConnection(app, c, cmdType)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: %v", err)
	}
	return app, err
}

func Run(app *fiber.App) {
	// Defina porta
	port := os.Getenv("APP_PORT")

	go func() {
		fmt.Printf("Servidor rodando na porta %s...\n", port)
		if err := app.Listen(":"+port); err != nil {
			log.Fatal("Erro ao rodar o servidor:", err)
			return
		}
	}()

	// Aguarde sinal de interrupção
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	if err := app.Shutdown(); 
	err != nil {
        log.Fatal("Erro ao encerrar servidor:", err)
		return 
    }

	//encerrando a fila
	queue.GetQueue().Close()
    queue.GetQueue().Wait()

	fmt.Println("Servidor encerrado!")
}
