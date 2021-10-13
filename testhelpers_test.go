package main

import (
	"encoding/json"
	"flag"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/sanity-io/litter"
)

var updateGoldenFiles = flag.Bool("updategolden", false, "set to update all golden files")
var removeAllGoldenFiles = flag.Bool("removeallgolden", false, "remove all golden files before updating them")

func ValidateWithGoldenFiles(t testing.TB, testResult interface{}, errCallback func()) {
	t.Helper()

	goldenFile := nameToFilepath(t.Name()) + ".golden"
	got := litter.Sdump(testResult) + "\n"

	if *updateGoldenFiles {
		if err := os.MkdirAll("testdata/", 0777); err != nil {
			t.Fatalf("cannot create testdate folder: %v", err)
		}

		if *removeAllGoldenFiles {
			if err := os.RemoveAll("testdate/*.golden"); err != nil {
				t.Fatalf("cannot remove files: %v", err)
			}
		}

		if err := os.WriteFile(goldenFile, []byte(got), 0644); err != nil {
			t.Fatalf("cannot write golden file: %v", err)
		}
	}

	wantRaw, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("cannot read goldenfile: %v", err)
	}

	want := string(wantRaw)

	if got == want {
		return
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(got),
		B:        difflib.SplitLines(want),
		FromFile: "Got",
		ToFile:   "Want",
		Context:  3,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		t.Fatalf("cannot create a unified diff: %s", err)
	}

	if errCallback != nil {
		errCallback()
	}

	t.Fatalf("Unexpected test result.\n%s", text)
}

func loadJSON(t *testing.T, name string) GameState {
	t.Helper()

	fpath := nameToFilepath(name) + ".json"
	data, err := os.ReadFile(fpath)
	if err != nil {
		t.Fatalf("reading file %q: %v", fpath, err)
	}

	state := GameState{}
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("unmarshaling data: %v", err)
	}

	return state
}

func nameToFilepath(testName string) string {
	path := regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(testName, "_")
	path = regexp.MustCompile(`_+`).ReplaceAllString(path, "_")
	path = strings.Trim(path, "_")
	return "testdata/" + path
}
