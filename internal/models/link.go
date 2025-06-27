package models

import "time"

type Link struct {
    ID        uint      `gorm:"primaryKey"`
    ShortCode string    `gorm:"uniqueIndex;size:6"`
    LongURL   string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Clicks    []Click
}

