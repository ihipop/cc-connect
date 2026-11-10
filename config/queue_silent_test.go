package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestQueueSilentConfig(t *testing.T) {
	for _, raw := range []string{"", "[queue]\nsilent = false", "[queue]\nsilent = true\nmax_depth = 2"} {
		var cfg Config
		if _, err := toml.Decode(raw, &cfg); err != nil {
			t.Fatal(err)
		}
		if want := raw == "[queue]\nsilent = true\nmax_depth = 2"; cfg.Queue.Silent != want {
			t.Fatalf("silent = %t, want %t for %q", cfg.Queue.Silent, want, raw)
		}
	}
}
