package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (v VideoModel) Insert(video *Video) error {
	query := `INSERT INTO videos (video_id, title, description, type, length, language, published_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING video_id, created_at, version`

	args := []any{video.VideoID, video.Title, video.Description, video.Type, video.Length, video.Language, video.PublishedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return v.DB.QueryRowContext(ctx, query, args...).Scan(&video.VideoID, &video.CreatedAt, &video.Version)
}

func (v VideoModel) Get(id string) (*Video, error) {
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

	err := v.DB.QueryRowContext(ctx, query, id).Scan(
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

func (v VideoModel) Update(video *Video) error {
	query := `
	UPDATE videos
	SET title = $1, description = $2, type = $3, length = $4, language = $5, published_at = $6, version = version + 1
	WHERE video_id = $8 AND version = $7
	RETURNING version`

	args := []any{
		video.Title,
		video.Description,
		video.Type,
		video.Length,
		video.Language,
		video.PublishedAt,
		video.Version,
		video.VideoID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := v.DB.QueryRowContext(ctx, query, args...).Scan(&video.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (v VideoModel) Delete(id string) error {
	if id == "" {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM videos
		WHERE video_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := v.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (v VideoModel) GetAll(title string, description string, filters Filters) ([]*Video, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(),  video_id, title, description, type, length, language, published_at, created_at, version
		FROM videos
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, description, filters.limit(), filters.offset()}

	rows, err := v.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	videos := []*Video{}

	for rows.Next() {
		var video Video

		err := rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}

		videos = append(videos, &video)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return videos, metadata, nil
}
