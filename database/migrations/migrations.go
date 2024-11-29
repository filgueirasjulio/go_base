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
        modelType := reflect.TypeOf(model)
        tableName := modelType.Name()
        
        // Verifica se a tabela existe
        if !m.DB.Migrator().HasTable(model) {
            m.DB.AutoMigrate(model)
            migrated = true
        } else {
            for i := 0; i < modelType.NumField(); i++ {
                field := modelType.Field(i)
                gormTag := field.Tag.Get("gorm")
                jsonTag := field.Tag.Get("json")
                
                // Ignora campos ignorados pelo GORM
                if gormTag == "-" || jsonTag == "-" {
                    continue
                }
                
                // Extrai nome da coluna
                columnName := helpers.ExtractJSONTag(jsonTag)
                if columnName == "" {
                    columnName = field.Name
                }
                
                // Verifica se a coluna existe
                if !m.DB.Migrator().HasColumn(tableName, columnName) {
                    m.DB.Migrator().AddColumn(tableName, columnName)
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
        models.RegisterHash{},
    }
}