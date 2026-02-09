package models

import (
	"time"
)

type Profile struct {
	UserID            string    `json:"user_id"`
	DisplayName       string    `json:"display_name"`
	Bio               string    `json:"bio"`
	SkillLevel        int       `json:"skill_level"`
	IsTemplateSeller  bool      `json:"is_template_seller"`
	RatingAvg         float64   `json:"rating_avg"`
	RatingCount       int       `json:"rating_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}