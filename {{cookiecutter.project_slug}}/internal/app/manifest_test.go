package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The service manifest follows the pulse standard: pulsekit never crashes the service, it skips a
// wrong item and records the reason in Problems() — here the list must be empty. Warnings() are
// texts over the standard limits (pulse will cut them). Register the diagnostic endpoints here as
// in app.go (with fake dependencies) and add a kit.CheckEndpoint call for each one — sample:
// internal/infra/pulsekit/example_test.go.
func TestManifest(t *testing.T) {
	kit := newPulsekit()
	require.Empty(t, kit.Problems(), "pulse standard violations: these items are not published")
	for _, w := range kit.Warnings() {
		t.Log("warning:", w)
	}

	m := kit.Manifest()

	require.Equal(t, 1, m.PulseManifest)
	assert.NotEmpty(t, m.Service.Owner.Team)
	assert.NotEmpty(t, m.Service.Description)
	assert.Contains(t, []string{"high", "medium", "low"}, m.Service.Criticality)
}
