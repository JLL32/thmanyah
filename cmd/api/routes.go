package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/videos", app.createVideoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/videos/:id", app.showVideoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/videos/:id", app.updateVideoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/videos/:id", app.deleteVideoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/videos", app.listVideosHandler)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return app.recoverPanic(router)
}
