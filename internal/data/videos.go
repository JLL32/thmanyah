package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JLL32/thmanyah/internal/validator"
)

type Video struct {
	VideoID     string    `json:"video_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Length      int       `json:"length"`
	Language    string    `json:"language"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int       `json:"version"`
}

func ValidateVideo(v *validator.Validator, video *Video) {
	v.Check(video.Title != "", "title", "must be provided")
	v.Check(len(video.Title) <= 500, "title", "must not be more than 100 bytes long")

	v.Check(video.Description != "", "description", "must be provided")
	v.Check(len(video.Description) <= 5000, "description", "must not be more than 5000 bytes long")

	v.Check(video.Type != "", "type", "must be provided")

	v.Check(video.Length > 0, "length", "must be greater than zero")

	v.Check(video.Language != "", "language", "must be provided")
	v.Check(len(video.Language) <= 2, "language", "must not be more than 50 bytes long")

	v.Check(video.PublishedAt.Before(time.Now()), "published_at", "must not be in the future")
}

type VideoModel struct {
	DB *sql.DB
}

func (m VideoModel) Insert(video *Video) error {
	query := `INSERT INTO videos (video_id, title, description, type, length, language, published_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING video_id, created_at, version`

	args := []any{video.VideoID, video.Title, video.Description, video.Type, video.Length, video.Language, video.PublishedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&video.VideoID, &video.CreatedAt, &video.Version)
}

func (m VideoModel) Get(id string) (*Video, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT video_id, title, description, type, length, language, published_at, created_at, version
	FROM videos
	WHERE video_id = $1
	`

	var video Video

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&video.VideoID,
		&video.Title,
		&video.Description,
		&video.Type,
		&video.Length,
		&video.Language,
		&video.PublishedAt,
		&video.CreatedAt,
		&video.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &video, nil
}
