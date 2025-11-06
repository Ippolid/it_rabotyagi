package models

import "time"

type Question struct {
    ID           int64     `json:"id"`
    Title        string    `json:"title"`
    Content      string    `json:"content"`
    Difficulty   *string   `json:"difficulty,omitempty"`
    Options      []string  `json:"options,omitempty"`
    CorrectAnswer *string  `json:"correctAnswer,omitempty"`
    Explanation  *string   `json:"explanation,omitempty"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
