package main

import (
	"context"
	"testing"
	"time"

	"github.com/amirzayi/clean_architect/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	cfg, err := config.LoadConfigOrDefault("")
	require.NoError(t, err)

	err = run(ctx, cfg)
	require.NoError(t, err)

	t.Cleanup(cancel)
}
