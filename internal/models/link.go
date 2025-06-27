package models

import "time"

type Link struct {
    ID         uint      `gorm:"primaryKey"`
    ShortCode  string    `gorm:"uniqueIndex;size:6"`
    LongURL    string    `gorm:"not null"`
    ClickCount uint      `gorm:"default:0"` // Compteur de clics pour optimiser les requêtes
    CreatedAt  time.Time
    UpdatedAt  time.Time
    Clicks     []Click   // Relation one-to-many avec Click
}

