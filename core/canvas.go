package core

import (
	"context"
	"time"
)

type (
	// CanvasMeta represents lightweight metadata for listing user-owned canvases.
	CanvasMeta struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		ThumbnailURL string    `json:"thumbnail,omitempty"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// Canvas represents the full metadata and content of a user-saved drawing.
	Canvas struct {
		ID           string    `json:"id"`
		UserID       string    `json:"-"` // Internal ownership field, never exposed in JSON.
		Name         string    `json:"name"`
		ThumbnailURL string    `json:"thumbnail,omitempty"`
		Data         []byte    `json:"data,omitempty"` // Full serialized canvas payload.
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// CanvasStore defines the persistence layer for user-owned canvases.
	CanvasStore interface {
		// List returns lightweight metadata for all canvases owned by a user.
		List(ctx context.Context, userID string) ([]*CanvasMeta, error)

		// Get returns a single canvas by its ID, ensuring it belongs to the user.
		Get(ctx context.Context, userID, id string) (*Canvas, error)

		// Save creates a new canvas or updates an existing one for the same user.
		Save(ctx context.Context, canvas *Canvas) error

		// Delete removes a canvas, ensuring it belongs to the user.
		Delete(ctx context.Context, userID, id string) error
	}
)