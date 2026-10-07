package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOnlyRootAndListedAdminsMaySetConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cc.yaml")
	if err := os.WriteFile(path, []byte("admins:\n  - admin-1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	prev := CONFIG_LOCATION
	CONFIG_LOCATION = path
	t.Cleanup(func() { CONFIG_LOCATION = prev })

	cases := map[string]bool{
		ROOT_ADMIN: true,
		"admin-1":  true,
		"stranger": false,
		"":         false,
	}
	for requestor, want := range cases {
		got, err := CanSetConfig(requestor)
		if err != nil {
			t.Fatalf("%q: %v", requestor, err)
		}
		if got != want {
			t.Errorf("%q: allowed=%v, want %v", requestor, got, want)
		}
	}
}

func TestARefusedRequestDoesNotWriteTheFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cc.yaml")
	before := []byte("admins:\n  - admin-1\n")
	if err := os.WriteFile(path, before, 0o600); err != nil {
		t.Fatal(err)
	}
	prev := CONFIG_LOCATION
	CONFIG_LOCATION = path
	t.Cleanup(func() { CONFIG_LOCATION = prev })

	// A stranger listing themselves in the request must not get in on it.
	defaults, _ := Config()
	defaults.Admins = append(defaults.Admins, "stranger")
	if _, err := SetConfig("stranger", defaults); err != ErrNotConfigAdmin {
		t.Fatalf("err = %v, want ErrNotConfigAdmin", err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(before) {
		t.Fatalf("config was rewritten: %q", after)
	}
}
