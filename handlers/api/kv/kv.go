package kv

import (
	"encoding/json"
	"excalidraw-complete/core"
	"excalidraw-complete/handlers/auth"
	"excalidraw-complete/middleware"
	"excalidraw-complete/stores"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type canvasPayload struct {
	Name      string          `json:"name"`
	Thumbnail string          `json:"thumbnail"`
	Data      json.RawMessage `json:"data"`
}

func HandleListCanvases(store stores.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*auth.AppClaims)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User claims not found"})
			return
		}

		canvases, err := store.List(r.Context(), claims.Subject)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userID": claims.Subject,
			}).Error("Failed to list canvases")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to list canvases"})
			return
		}

		if canvases == nil {
			canvases = []*core.Canvas{}
		}

		render.JSON(w, r, canvases)
	}
}

func HandleGetCanvas(store stores.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*auth.AppClaims)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User claims not found"})
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Canvas key is required"})
			return
		}

		canvas, err := store.Get(r.Context(), claims.Subject, key)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userID": claims.Subject,
				"key":    key,
			}).Warn("Failed to get canvas")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Canvas not found"})
			return
		}

		render.JSON(w, r, map[string]any{
			"id":        canvas.ID,
			"name":      canvas.Name,
			"thumbnail": canvas.Thumbnail,
			"data":      json.RawMessage(canvas.Data),
			"createdAt": canvas.CreatedAt,
			"updatedAt": canvas.UpdatedAt,
		})
	}
}

func HandleSaveCanvas(store stores.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*auth.AppClaims)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User claims not found"})
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Canvas key is required"})
			return
		}

		defer r.Body.Close()

		var payload canvasPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"key":   key,
			}).Error("Failed to decode request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid request body"})
			return
		}

		if payload.Name == "" {
			payload.Name = key
		}

		canvas := &core.Canvas{
			ID:        key,
			UserID:    claims.Subject,
			Name:      payload.Name,
			Thumbnail: payload.Thumbnail,
			Data:      []byte(payload.Data),
		}

		if err := store.Save(r.Context(), canvas); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userID": claims.Subject,
				"key":    key,
			}).Error("Failed to save canvas")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to save canvas"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]any{
			"id":        canvas.ID,
			"name":      canvas.Name,
			"thumbnail": canvas.Thumbnail,
		})
	}
}

func HandleDeleteCanvas(store stores.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*auth.AppClaims)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "User claims not found"})
			return
		}

		key := chi.URLParam(r, "key")
		if key == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Canvas key is required"})
			return
		}

		if err := store.Delete(r.Context(), claims.Subject, key); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"userID": claims.Subject,
				"key":    key,
			}).Error("Failed to delete canvas")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete canvas"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"status": "deleted"})
	}
}