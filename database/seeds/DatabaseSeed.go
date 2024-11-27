package seeds

import (
	"log"
	"gorm.io/gorm"
)

type Seeder struct {
	DB *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{DB: db}
}

// Função que executa todas as seeders
func (s *Seeder) RunAllSeeds() {
	s.RunModelSeed("User")

	log.Println("Seeds executadas com sucesso")
}

// Função que executa a seeder de um único modelo
func (s *Seeder) RunModelSeed(model string) {
	switch model {
	case "User":
		SeedUser(s.DB)
	default:
		log.Printf("Modelo %s não encontrado ou não suportado", model)
	}

	log.Printf("%sSeed executada com sucesso", model)
}