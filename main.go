package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	queue "base/app"
	"base/connections"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	zerolog "github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var app *fiber.App
var ctx *cli.Context
var zeroLog zerolog.Logger


func init() {
	app = fiber.New()

	//fiber middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	//salvar log em txt
	zeroLog = zerolog.New(os.Stdout).With().Caller().Logger()

	initLogger()

	// Carrega as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	//iniciar fila
	capacity, _ := strconv.Atoi(os.Getenv("QUEUE_CAPACITY"))

	queue.InitQueue(capacity)

	//conectar ao banco
	Teste(app, ctx, "test", zeroLog);

	port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000" 
    }

    log.Printf("Servidor rodando na porta %s...", port)

    if err := app.Listen(":" + port); err != nil {
        log.Fatalf("Erro ao iniciar servidor: %v", err)
    }
}

func main() {
	cmd := &cli.App{
		Name:        "base",
		Usage:       "Versão de Desenvolvimento",
		Version:     "1.0.0",
		UsageText:   "Rode main [global options] command [command options] [arguments...]",
		Description: "Base para desenvolvimento em Go",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Executando em modo padrão",
				Action: func(ctx *cli.Context) error {
					// Conectar ao banco e rodar o serviço
					app, err := Teste(app, ctx, "run", zeroLog)
					if err != nil {
						log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
					}

					Run(app)
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Executando teste de conexão",
				Action: func(ctx *cli.Context) error {
					_, err := Teste(app, ctx, "test", zeroLog)
					if err != nil {
						log.Fatalf("Falha no teste de conexão: %v", err)
					}

					log.Println("Conexão com o banco de dados foi bem-sucedida!")
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "Executar as migrations",
				Action: func(ctx *cli.Context) error {
					_, err := Teste(app, ctx, "migrate", zeroLog)
					if err != nil {
						log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
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
					_, err := Teste(app, ctx, "seed", zeroLog)
					if err != nil {
						log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
					}
					return nil
				},
			},
		},
		Compiled:  time.Now(),
		Authors:   []*cli.Author{},
		Copyright: fmt.Sprintf("© %d Go Base - todos direitos reservados", time.Now().Year()),
	}

	cmd.Run(os.Args)
}

func Teste(app *fiber.App, c *cli.Context, cmdType string, logger zerolog.Logger) (*fiber.App, error) {
	zeroLog.Info().Msg("Teste de conexão iniciado")
	startTime := time.Now()

	err := connections.GetDatabaseConnection(app, c, cmdType, logger)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	elapsedTime := time.Since(startTime)
	zeroLog.Info().Msg("✅ Finalizado o teste de conexão") 
	zeroLog.Info().Msgf("Tempo de execução: %s",  elapsedTime)

	return app, err
}

func Run(app *fiber.App) {

	zeroLog.Info().Msgf("Servidor iniciado")
	startTime := time.Now()

	// Porta
	port := os.Getenv("APP_PORT")

	go func() {
		fmt.Printf("Servidor rodando na porta %s...\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Erro ao rodar o servidor: %v", err)
			return
		}
	}()

	// Aguarde sinal de interrupção
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
		return
	}

	//encerrando a fila
	queue.GetQueue().Close()
	queue.GetQueue().Wait()
		
	elapsedTime := time.Since(startTime)
	zeroLog.Info().Msg("✅ Servidor finalizado") 
	zeroLog.Info().Msgf("Tempo de execução: %s",  elapsedTime)
}

func initLogger() {
    logsDir := "logs"
    if _, err := os.Stat(logsDir); os.IsNotExist(err) {
        os.Mkdir(logsDir, 0755)
    }

    filename := filepath.Join(logsDir, "log_"+time.Now().Format("2006-01-02")+".txt")
    f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Printf("Erro ao criar arquivo de log: %v", err)
        os.Exit(1)
    }

    // Configurar o formato da hora
    zerolog.TimeFieldFormat = "15:04:05"

    zeroLog = zerolog.New(io.MultiWriter(zerolog.NewConsoleWriter(), f)).
        With().
        Caller().
        Timestamp().
        Logger()

    zeroLog.Info().Msg("Logger iniciado")
}
