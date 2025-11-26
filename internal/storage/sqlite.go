package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS attestations (
		id TEXT PRIMARY KEY,
		tool_name TEXT NOT NULL,
		input TEXT NOT NULL,
		output TEXT NOT NULL,
		origin TEXT NOT NULL,
		confidence REAL NOT NULL,
		explanation TEXT NOT NULL,
		request_id TEXT,
		created_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_attestations_created_at ON attestations(created_at);
	CREATE INDEX IF NOT EXISTS idx_attestations_origin ON attestations(origin);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) SaveAttestation(ctx context.Context, a *Attestation) error {
	inputJSON, err := json.Marshal(a.Input)
	if err != nil {
		return err
	}
	outputJSON, err := json.Marshal(a.Output)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO attestations (id, tool_name, input, output, origin, confidence, explanation, request_id, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.ExecContext(ctx, query,
		a.ID, a.ToolName, string(inputJSON), string(outputJSON),
		a.Origin, a.Confidence, a.Explanation, a.RequestID, a.CreatedAt,
	)
	return err
}

func (s *SQLiteStore) GetAttestation(ctx context.Context, id string) (*Attestation, error) {
	query := `SELECT id, tool_name, input, output, origin, confidence, explanation, request_id, created_at FROM attestations WHERE id = ?`
	row := s.db.QueryRowContext(ctx, query, id)

	var a Attestation
	var inputJSON, outputJSON string
	err := row.Scan(&a.ID, &a.ToolName, &inputJSON, &outputJSON, &a.Origin, &a.Confidence, &a.Explanation, &a.RequestID, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(inputJSON), &a.Input)
	json.Unmarshal([]byte(outputJSON), &a.Output)

	return &a, nil
}

func (s *SQLiteStore) ListAttestations(ctx context.Context, limit, offset int) ([]*Attestation, error) {
	query := `SELECT id, tool_name, input, output, origin, confidence, explanation, request_id, created_at 
	          FROM attestations ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attestations []*Attestation
	for rows.Next() {
		var a Attestation
		var inputJSON, outputJSON string
		if err := rows.Scan(&a.ID, &a.ToolName, &inputJSON, &outputJSON, &a.Origin, &a.Confidence, &a.Explanation, &a.RequestID, &a.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(inputJSON), &a.Input)
		json.Unmarshal([]byte(outputJSON), &a.Output)
		attestations = append(attestations, &a)
	}

	return attestations, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

type MemoryStore struct {
	attestations map[string]*Attestation
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{attestations: make(map[string]*Attestation)}
}

func (s *MemoryStore) SaveAttestation(ctx context.Context, a *Attestation) error {
	s.attestations[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAttestation(ctx context.Context, id string) (*Attestation, error) {
	return s.attestations[id], nil
}

func (s *MemoryStore) ListAttestations(ctx context.Context, limit, offset int) ([]*Attestation, error) {
	result := make([]*Attestation, 0, len(s.attestations))
	for _, a := range s.attestations {
		result = append(result, a)
	}
	if offset >= len(result) {
		return []*Attestation{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (s *MemoryStore) Close() error {
	return nil
}

func GenerateID() string {
	return time.Now().Format("20060102150405") + "-" + randomHex(4)
}

func randomHex(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = "0123456789abcdef"[time.Now().UnixNano()%16]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
