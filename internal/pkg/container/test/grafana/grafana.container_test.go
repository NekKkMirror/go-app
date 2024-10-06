package grafanacontainer

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StartGrafanaContainer(t *testing.T) {
	ctx := context.Background()
	containerInfo, err := Start(ctx, t)
	require.NoError(t, err)

	assert.NotNil(t, containerInfo)
	assert.Equal(t, "admin", containerInfo.AdminUser)

	url := fmt.Sprintf("http://%s:%d/login", containerInfo.Host, containerInfo.Port)
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
