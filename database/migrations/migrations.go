package migrations

import (
	"log"
	"gorm.io/gorm"
	"reflect"

	"base/app/models"
	"base/utils/helpers"
)

type Migration struct {
	DB *gorm.DB
}

func NewMigration(db *gorm.DB) *Migration {
	return &Migration{DB: db}
}

// RunMigrations executa todas as migrações dos models
func (m *Migration) RunMigrations() {
	models := getModels()
    migrated := m.checkAndMigrateTablesAndColumns(models)

    if !migrated {
        log.Println("Nada para migrar.")
    }
}

// checkAndMigrateTableAndColumns verifica se a tabela e as colunas de um modelo existem no banco de dados e faz as migrações necessárias
func (m *Migration) checkAndMigrateTablesAndColumns(models []interface{}) bool {
    migrated := false
    for _, model := range models {
		if model == nil {
            log.Println("Modelo nulo")
            continue
        }

        if m.DB == nil {
            log.Fatal("Banco de dados não inicializado")
            return false
        }

        modelType := reflect.TypeOf(model) // Obtém o tipo do struct

        // Verifica se a tabela existe
        tableExists := m.DB.Migrator().HasTable(model)
        if !tableExists {
            log.Printf("Criando tabela para o modelo '%s'.", modelType.Name())
            m.DB.AutoMigrate(model) // Cria a tabela
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
                if !m.DB.Migrator().HasColumn(model, columnName) && columnName != "validation_code" {
                    log.Printf("Adicionando coluna '%s' à tabela '%s'.", columnName, modelType.Name())
                    m.DB.Migrator().AddColumn(model, columnName)
                    migrated = true
                }
            }
        }
    }

    return migrated
}

func getModels() []interface{} {
    return []interface{}{
        models.User{},
        models.ValidationCode{},
    }
}