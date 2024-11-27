package migrations

import (
	"log"
	"gorm.io/gorm"
	"reflect"

	"tradeapi/app/models"
	"tradeapi/utils/helpers"
)

type Migration struct {
	DB *gorm.DB
}

func NewMigration(db *gorm.DB) *Migration {
	return &Migration{DB: db}
}

// RunMigrations executa todas as migrações dos models
func (m *Migration) RunMigrations() {

	migrated := checkAndMigrateTableAndColumns(m.DB, &models.User{}) 

	// Verifica se não houve migração executada
	if !migrated {
		log.Println("Nada para migrar.")
	}
}

// checkAndMigrateTableAndColumns verifica se a tabela e as colunas de um modelo existem no banco de dados e faz as migrações necessárias
func checkAndMigrateTableAndColumns(db *gorm.DB, model interface{}) bool {
	migrated := false
	modelType := reflect.TypeOf(model).Elem() // Obtém o tipo do struct 

	// Verifica se a tabela existe
	tableExists := db.Migrator().HasTable(model)
	if !tableExists {
		log.Printf("Criando tabela para o modelo '%s'.", modelType.Name())
		db.AutoMigrate(model) // Cria a tabela
		migrated = true
	} else {
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			jsonTag := field.Tag.Get("json")
			columnName := helpers.ExtractJSONTag(jsonTag)

			// Ignorar campos sem a tag json ou que estão marcados para ignorar (json:"-")
			if columnName == "" || columnName == "-" {
				continue
			}

			// Verifica se a coluna existe
			if !db.Migrator().HasColumn(model, columnName) {
				log.Printf("Adicionando coluna '%s' à tabela '%s'.", columnName, modelType.Name())
				db.Migrator().AddColumn(model, columnName) 
				migrated = true
			}
		}
	}

	return migrated
}

