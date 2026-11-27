package cmd

import "testing"

func TestSearchListCommandsRejectPositionalArgs(t *testing.T) {
	root := newRootCmd()
	for _, path := range [][]string{
		{"dashboard", "list"},
		{"user", "list"},
		{"team", "list"},
		{"service-account", "list"},
		{"playlist", "list"},
		{"library-element", "list"},
	} {
		cmd, _, err := root.Find(path)
		if err != nil {
			t.Fatalf("%v: %v", path, err)
		}
		if err := cmd.ValidateArgs([]string{"production"}); err == nil {
			t.Errorf("%v accepted a stray positional argument", path)
		}
	}
}
