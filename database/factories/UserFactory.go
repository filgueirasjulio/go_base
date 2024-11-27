package factories

import (
	"tradeapi/app/models"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"gorm.io/gorm"
	"time"
)

// FactoryUser cria um ou mais usuários de maneira flexível.
func FactoryUser(db *gorm.DB, quantity int) ([]models.User, error) {
	users := []models.User{}

	for i := 0; i < quantity; i++ {
		user := models.User{
			Name:  faker.Name(),
			Email: faker.Email(),
			CreatedAt: time.Now(),
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Erro ao gerar a senha para o usuário: %v", err)
			return nil, err
		}
		user.Password = string(hashedPassword)

		users = append(users, user)
	}

	if err := db.Create(&users).Error; err != nil {
		log.Printf("Erro ao criar usuários: %v", err)
		return nil, err
	}

	return users, nil
}
