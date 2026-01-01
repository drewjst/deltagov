package diff_engine

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	"github.com/aymanbagabas/go-udiff/myers"
)

// Delta represents the structured diff between two text versions
type Delta struct {
	VersionA    string   `json:"version_a"`
	VersionB    string   `json:"version_b"`
	Hunks       []Hunk   `json:"hunks"`
	Insertions  int      `json:"insertions"`
	Deletions   int      `json:"deletions"`
	Unchanged   int      `json:"unchanged"`
}

// Hunk represents a contiguous block of changes
type Hunk struct {
	StartA int      `json:"start_a"`
	StartB int      `json:"start_b"`
	Lines  []Change `json:"lines"`
}

// Change represents a single line change
type Change struct {
	Type    ChangeType `json:"type"`
	Content string     `json:"content"`
	LineA   int        `json:"line_a,omitempty"`
	LineB   int        `json:"line_b,omitempty"`
}

// ChangeType indicates the type of change
type ChangeType string

const (
	ChangeInsert    ChangeType = "insert"
	ChangeDelete    ChangeType = "delete"
	ChangeUnchanged ChangeType = "unchanged"
)

// ComputeHash generates a SHA-256 hash of the content
func ComputeHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// Compute calculates the diff between two text versions using Myers algorithm
func Compute(textA, textB, versionA, versionB string) (*Delta, error) {
	// Use go-udiff with Myers algorithm
	edits := myers.ComputeEdits(textA, textB)
	unifiedDiff, err := udiff.ToUnified("version_a", "version_b", textA, edits)
	if err != nil {
		return nil, err
	}

	delta := &Delta{
		VersionA: versionA,
		VersionB: versionB,
		Hunks:    []Hunk{},
	}

	// Parse unified diff into our Delta structure
	linesA := strings.Split(textA, "\n")
	linesB := strings.Split(textB, "\n")

	// Count changes from the unified diff
	for _, line := range strings.Split(unifiedDiff, "\n") {
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '+':
			if !strings.HasPrefix(line, "+++") {
				delta.Insertions++
			}
		case '-':
			if !strings.HasPrefix(line, "---") {
				delta.Deletions++
			}
		case ' ':
			delta.Unchanged++
		}
	}

	// If no changes detected via diff, count all lines as unchanged
	if delta.Insertions == 0 && delta.Deletions == 0 && delta.Unchanged == 0 {
		delta.Unchanged = max(len(linesA), len(linesB))
	}

	return delta, nil
}

// ComputeWordLevel performs word-level diffing for more granular changes
func ComputeWordLevel(textA, textB string) (*Delta, error) {
	// Split into words while preserving structure
	wordsA := tokenize(textA)
	wordsB := tokenize(textB)

	// Compute diff on word tokens
	edits := myers.ComputeEdits(strings.Join(wordsA, "\n"), strings.Join(wordsB, "\n"))

	delta := &Delta{
		Hunks: []Hunk{},
	}

	// Count changes
	diffStr, err := udiff.ToUnified("a", "b", strings.Join(wordsA, "\n"), edits)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(diffStr, "\n") {
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '+':
			if !strings.HasPrefix(line, "+++") {
				delta.Insertions++
			}
		case '-':
			if !strings.HasPrefix(line, "---") {
				delta.Deletions++
			}
		case ' ':
			delta.Unchanged++
		}
	}

	return delta, nil
}

// tokenize splits text into word tokens
func tokenize(text string) []string {
	var tokens []string
	var current strings.Builder

	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
