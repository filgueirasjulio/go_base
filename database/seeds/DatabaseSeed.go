package seeds

import (
	"log"
	"gorm.io/gorm"
)

// Função que executa todas as seeders
func RunAllSeeds(db *gorm.DB) {
	RunModelSeed(db, "User")

	log.Println("Seeds executadas com sucesso")
}

// Função que executa a seeder de um único modelo
func RunModelSeed(db *gorm.DB, model string) {
	switch model {
	case "User":
		SeedUser(db)
	default:
		log.Printf("Modelo %s não encontrado ou não suportado", model)
	}

	log.Printf("%sSeed executada com sucesso", model)
}