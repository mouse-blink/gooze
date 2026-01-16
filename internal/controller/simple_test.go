package controller

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/spf13/cobra"
)

func TestSimpleUI_DisplayEstimation_PrintsTable(t *testing.T) {
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	ui := NewSimpleUI(cmd)

	mutations := []m.Mutation{
		{Source: m.Source{Origin: &m.File{Path: "path/a.go"}}},
		{Source: m.Source{Origin: &m.File{Path: "path/a.go"}}},
		{Source: m.Source{Origin: &m.File{Path: "path/b.go"}}},
	}

	if err := ui.DisplayEstimation(mutations, nil); err != nil {
		t.Fatalf("DisplayEstimation() error = %v", err)
	}

	output := buf.String()

	for _, want := range []string{
		"path/a.go",
		"path/b.go",
		"2",
		"1",
		"TOTAL FILES 2",
		"3",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestSimpleUI_DisplayEstimation_Error(t *testing.T) {
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	ui := NewSimpleUI(cmd)
	boom := errors.New("boom")

	if err := ui.DisplayEstimation(nil, boom); err == nil {
		t.Fatalf("DisplayEstimation() expected error")
	}

	output := buf.String()
	if !strings.Contains(output, "estimation error: boom") {
		t.Fatalf("output missing error message\noutput:\n%s", output)
	}
}
