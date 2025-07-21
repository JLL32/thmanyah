package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JLL32/thmanyah/internal/data"
	"github.com/JLL32/thmanyah/internal/validator"
)

func (app *application) createVideoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		VideoID     string    `json:"video_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Type        string    `json:"type"`
		Language    string    `json:"language"`
		Length      int       `json:"length"`
		PublishedAt time.Time `json:"published_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	video := &data.Video{
		VideoID:     input.VideoID,
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Language:    input.Language,
		Length:      input.Length,
		PublishedAt: input.PublishedAt,
	}

	v := validator.New()
	if data.ValidateVideo(v, video); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Videos.Insert(video)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/videos/%s", video.VideoID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"video": video}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showVideoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	video, err := app.models.Videos.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"video": video}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the updatevideoHandler function
}

func (app *application) deleteVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the deletevideoHandler function
}

func (app *application) listVideosHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the listVideosHandler function
}
