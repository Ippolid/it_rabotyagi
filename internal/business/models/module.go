package models

import "time"

type Module struct {
    ID          int64     `json:"id"`
    CourseID    int64     `json:"courseId"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    ModuleOrder int       `json:"moduleOrder"`
    CreatedAt   time.Time `json:"createdAt"`
    EditedAt    time.Time `json:"editedAt"`
}
