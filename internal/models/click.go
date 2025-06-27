package models

import "time"

type Click struct {
    ID        uint      `gorm:"primaryKey"`
    LinkID    uint      `gorm:"index"` // clé étrangère vers Link
    Timestamp time.Time 
    UserAgent string
    IP        string
}
