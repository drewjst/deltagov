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
	// 3 lines of context around changes (standard unified diff format)
	unifiedDiff, err := udiff.ToUnified("version_a", "version_b", textA, edits, 3)
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
	// Split by lines for line-level diffing with word context
	linesA := strings.Split(textA, "\n")
	linesB := strings.Split(textB, "\n")

	// Compute diff on lines
	edits := myers.ComputeEdits(textA, textB)

	delta := &Delta{
		Hunks: []Hunk{},
	}

	// Generate unified diff
	diffStr, err := udiff.ToUnified("a", "b", textA, edits, 3)
	if err != nil {
		return nil, err
	}

	// Parse the unified diff to extract changes
	var currentHunk *Hunk
	lineNumA := 1
	lineNumB := 1

	for _, line := range strings.Split(diffStr, "\n") {
		if len(line) == 0 {
			continue
		}

		// Skip header lines
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			continue
		}

		// Parse hunk header
		if strings.HasPrefix(line, "@@") {
			if currentHunk != nil {
				delta.Hunks = append(delta.Hunks, *currentHunk)
			}
			currentHunk = &Hunk{
				StartA: lineNumA,
				StartB: lineNumB,
				Lines:  []Change{},
			}
			continue
		}

		if currentHunk == nil {
			continue
		}

		switch line[0] {
		case '+':
			content := ""
			if len(line) > 1 {
				content = line[1:]
			}
			currentHunk.Lines = append(currentHunk.Lines, Change{
				Type:    ChangeInsert,
				Content: content,
				LineB:   lineNumB,
			})
			delta.Insertions++
			lineNumB++
		case '-':
			content := ""
			if len(line) > 1 {
				content = line[1:]
			}
			currentHunk.Lines = append(currentHunk.Lines, Change{
				Type:    ChangeDelete,
				Content: content,
				LineA:   lineNumA,
			})
			delta.Deletions++
			lineNumA++
		case ' ':
			content := ""
			if len(line) > 1 {
				content = line[1:]
			}
			currentHunk.Lines = append(currentHunk.Lines, Change{
				Type:    ChangeUnchanged,
				Content: content,
				LineA:   lineNumA,
				LineB:   lineNumB,
			})
			delta.Unchanged++
			lineNumA++
			lineNumB++
		}
	}

	// Add the last hunk
	if currentHunk != nil && len(currentHunk.Lines) > 0 {
		delta.Hunks = append(delta.Hunks, *currentHunk)
	}

	// If no changes detected, add all lines as unchanged
	if len(delta.Hunks) == 0 {
		hunk := Hunk{StartA: 1, StartB: 1, Lines: []Change{}}
		maxLines := max(len(linesA), len(linesB))
		for i := 0; i < maxLines; i++ {
			content := ""
			if i < len(linesA) {
				content = linesA[i]
			} else if i < len(linesB) {
				content = linesB[i]
			}
			hunk.Lines = append(hunk.Lines, Change{
				Type:    ChangeUnchanged,
				Content: content,
				LineA:   i + 1,
				LineB:   i + 1,
			})
			delta.Unchanged++
		}
		delta.Hunks = append(delta.Hunks, hunk)
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
