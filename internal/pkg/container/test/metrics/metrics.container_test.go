package metricscontainer

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_StartMetricsContainer launches the metrics container and checks the availability of metrics.
func Test_StartMetricsContainer(t *testing.T) {
	ctx := context.Background()
	containerInfo, err := StartContainer(ctx, t)
	require.NoError(t, err)

	assert.NotNil(t, containerInfo)

	url := fmt.Sprintf("http://%s:%d/metrics", containerInfo.Host, containerInfo.Port)
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
