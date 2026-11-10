package atlas

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Corrections are fixes for known errors in the source data. They are applied
// at load time so the cleaned JSON stays untouched and every fix is documented.
type personFix struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Birth json.RawMessage `json:"birth"` // absent: leave, null: clear, number: set
	Death json.RawMessage `json:"death"`
	Note  string          `json:"note"`
}

type eventFix struct {
	ID        string `json:"id"`
	StartDate string `json:"startDate"`
	Duration  string `json:"duration"`
	Note      string `json:"note"`
}

// artFix re-points (or, with an empty event, un-matches) an engraving whose automatic match was wrong.
type artFix struct {
	Title string `json:"title"` // case-insensitive part of the picture's title
	Event string `json:"event"`
	Note  string `json:"note"`
}

type corrections struct {
	People []personFix `json:"people"`
	Events []eventFix  `json:"events"`
	Art    []artFix    `json:"art"`
}

func yearFix(raw json.RawMessage, cur *int) (*int, error) {
	if len(raw) == 0 {
		return cur, nil
	}
	if string(raw) == "null" {
		return nil, nil
	}
	var y int
	if err := json.Unmarshal(raw, &y); err != nil {
		return nil, err
	}
	return &y, nil
}

func (a *Atlas) applyCorrections(file string) error {
	b, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var c corrections
	if err := json.Unmarshal(b, &c); err != nil {
		return fmt.Errorf("corrections: %w", err)
	}
	for _, f := range c.People {
		p := a.PersonByID[f.ID]
		if p == nil {
			return fmt.Errorf("corrections: unknown person %s (%s)", f.ID, f.Name)
		}
		if p.Born, err = yearFix(f.Birth, p.Born); err != nil {
			return fmt.Errorf("corrections: %s: %w", f.Name, err)
		}
		if p.Died, err = yearFix(f.Death, p.Died); err != nil {
			return fmt.Errorf("corrections: %s: %w", f.Name, err)
		}
	}
	a.artFixes = c.Art
	for _, f := range c.Events {
		e := a.EventByID[f.ID]
		if e == nil {
			return fmt.Errorf("corrections: unknown event %s", f.ID)
		}
		if f.StartDate != "" {
			e.Year, e.Date = parseYear(f.StartDate)
		}
		if f.Duration != "" {
			e.Days = parseDays(f.Duration)
		}
	}
	return nil
}
