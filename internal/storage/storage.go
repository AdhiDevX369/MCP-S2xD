package storage

import (
	"context"
	"time"
)

type Attestation struct {
	ID          string         `json:"id"`
	ToolName    string         `json:"tool_name"`
	Input       map[string]any `json:"input"`
	Output      map[string]any `json:"output"`
	Origin      string         `json:"origin"`
	Confidence  float64        `json:"confidence"`
	Explanation string         `json:"explanation"`
	RequestID   string         `json:"request_id"`
	CreatedAt   time.Time      `json:"created_at"`
}

type Store interface {
	SaveAttestation(ctx context.Context, a *Attestation) error
	GetAttestation(ctx context.Context, id string) (*Attestation, error)
	ListAttestations(ctx context.Context, limit, offset int) ([]*Attestation, error)
	Close() error
}
