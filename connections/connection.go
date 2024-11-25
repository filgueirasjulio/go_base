package connections

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	"log"
	"os"
)

// GetDatabaseConnection cria e retorna uma conexão com o banco de dados usando o GORM
func GetDatabaseConnection() (*gorm.DB, error) {
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
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	// Testar a conexão (sem Ping)
	sqlDB, err := db.DB() 
	if err != nil {
		return nil, fmt.Errorf("erro ao obter DB do GORM: %v", err)
	}

	// Testa a conexão com o banco
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão com o banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados foi bem-sucedida!")
	return db, nil
}

