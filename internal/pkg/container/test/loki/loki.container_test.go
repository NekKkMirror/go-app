package lokicontainer

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StartLokiContainer(t *testing.T) {
	ctx := context.Background()
	lokiURL, err := Start(ctx, t)

	if assert.NoError(t, err) {
		readyURL := lokiURL + "/ready"
		resp, err := http.Get(readyURL)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	}
}
