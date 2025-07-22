package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

	if version := r.Header.Get("X-Expected-Version"); version != "" {
		if strconv.Itoa(int(video.Version)) != version {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Title *string `json:"title"`
		Description *string `json:"description"`
		Type *string `json:"type"`
		Length *int `json:"length"`
		Language *string `json:"language"`
		PublishedAt *time.Time `json:"published_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		video.Title = *input.Title
	}
	if input.Description != nil {
		video.Description = *input.Description
	}
	if input.Type != nil {
		video.Type = *input.Type
	}
	if input.Length != nil {
		video.Length = *input.Length
	}
	if input.Language != nil {
		video.Language = *input.Language
	}
	if input.PublishedAt != nil {
		video.PublishedAt = *input.PublishedAt
	}

	v := validator.New()
	if data.ValidateVideo(v, video); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Videos.Update(video)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"video": video}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteVideoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Videos.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "video successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listVideosHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Description string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Description =  app.readString(qs, "description", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"video_id", "title", "description", "length", "type", "-video_id", "-title", "-description", "-length", "-type"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	videos, metadata, err := app.models.Videos.GetAll(input.Title, input.Description, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "videos": videos}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
