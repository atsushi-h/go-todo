package database

import (
    "log"
    
    "go-todo/model"
    "gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
    log.Println("Starting database migration...")
    
    if err := db.AutoMigrate(
        &model.Todo{},
    ); err != nil {
        return err
    }
    
    log.Println("Database migration completed successfully.")
    return nil
}

func DropAllTables(db *gorm.DB) error {
    log.Println("Dropping all tables...")
    
    if err := db.Migrator().DropTable(
        &model.Todo{},
    ); err != nil {
        return err
    }
    
    log.Println("All tables dropped successfully.")
    return nil
}

func ResetDatabase(db *gorm.DB) error {
    if err := DropAllTables(db); err != nil {
        return err
    }
    return Migrate(db)
}
