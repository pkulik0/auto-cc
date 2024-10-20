package srt

import (
	"fmt"
	"strconv"
	"strings"
)

type line struct {
	Number int
	Start  string
	End    string
	Text   string
}

// Srt represents subtitles in the SRT format.
type Srt struct {
	Lines []line
}

// New creates a new SRT instance.
func Parse(data string) (*Srt, error) {
	lines := strings.Split(data, "\n")

	getLine := func(i int) (string, error) {
		if i < 0 || i >= len(lines) {
			return "", fmt.Errorf("line %d out of range", i)
		}
		return strings.Trim(lines[i], " \t\n"), nil
	}

	srt := &Srt{}
	for i := 0; i < len(lines); i++ {
		line := line{}

		// Parse the line number.
		lineText, err := getLine(i)
		if err != nil {
			return nil, err
		}
		if lineText == "" {
			continue
		}

		number, err := strconv.Atoi(lineText)
		if err != nil {
			return nil, err
		}
		line.Number = number

		// Parse the time range.
		i++

		lineText, err = getLine(i)
		if err != nil {
			return nil, err
		}

		timeRange := strings.Split(strings.TrimSpace(lines[i]), " --> ")
		if len(timeRange) != 2 {
			return nil, fmt.Errorf("invalid time range: %s", lineText)
		}
		line.Start = timeRange[0]
		line.End = timeRange[1]

		// Parse the text.
		i++
		lineText, err = getLine(i)
		if err != nil {
			return nil, err
		}
		line.Text = lineText

		srt.Lines = append(srt.Lines, line)
	}

	return srt, nil
}

// String serializes the SRT instance to a string.
func (s *Srt) String() string {
	var sb strings.Builder
	for _, l := range s.Lines {
		sb.WriteString(strconv.Itoa(l.Number))
		sb.WriteString("\n")
		sb.WriteString(l.Start)
		sb.WriteString(" --> ")
		sb.WriteString(l.End)
		sb.WriteString("\n")
		sb.WriteString(l.Text)
		sb.WriteString("\n\n")
	}
	return sb.String()
}

// Text returns the text of the subtitles.
func (s *Srt) Text() []string {
	var text []string
	for _, l := range s.Lines {
		text = append(text, l.Text)
	}
	return text
}

// ReplaceText replaces the text of the subtitles.
func (s *Srt) ReplaceText(text []string) error {
	if len(text) != len(s.Lines) {
		return fmt.Errorf("invalid text length: %d, expected: %d", len(text), len(s.Lines))
	}

	for i := range s.Lines {
		s.Lines[i].Text = text[i]
	}

	return nil
}

// Clone creates a deep copy of the SRT instance.
func (s *Srt) Clone() *Srt {
	clone := &Srt{}
	for _, l := range s.Lines {
		clone.Lines = append(clone.Lines, l)
	}
	return clone
}
