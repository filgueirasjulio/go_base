package seeds

import (
	"log"
	"tradeapi/database/factories"
	"gorm.io/gorm"
)

// Função que popula a tabela "users" com 5 usuários fictícios
func SeedUser(db *gorm.DB) {
	_, err := factories.FactoryUser(db, 5)
	if err != nil {
		log.Printf("Erro ao criar usuários: %v", err)
		return
	}
}
