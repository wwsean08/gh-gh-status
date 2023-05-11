package status

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifyConstants(t *testing.T) {
	require.Equal(t, "operational", COMPONENT_OPERATIONAL)
	require.Equal(t, "degraded_performance", COMPONENT_DEGREDADED_PERFORMANCE)
	require.Equal(t, "partial_outage", COMPONENT_PARTIAL_OUTAGE)
	require.Equal(t, "major_outage", COMPONENT_MAJOR_OUTAGE)
}
