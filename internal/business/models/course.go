package models

import "time"

type Course struct {
    ID          int64
    Title       string
    Description string
    IsPublished bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}


