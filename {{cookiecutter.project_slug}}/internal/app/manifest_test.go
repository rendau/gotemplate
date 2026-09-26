package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The service manifest builds without panics: pulsekit checks the pulse standard at
// registration (text limits, secrets, parameters), so a mistake fails here, not after deploy.
// For every diagnostic endpoint add a kit.CheckEndpoint call with fake dependencies —
// sample: internal/infra/pulsekit/example_test.go.
func TestManifest(t *testing.T) {
	var m = newPulsekit().Manifest()

	require.Equal(t, 1, m.PulseManifest)
	assert.NotEmpty(t, m.Service.Owner.Team)
	assert.NotEmpty(t, m.Service.Description)
	assert.Contains(t, []string{"high", "medium", "low"}, m.Service.Criticality)
}
