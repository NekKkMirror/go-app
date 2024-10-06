package grafana

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	grafanacontainer "github.com/NekKkMirror/go-app/internal/pkg/container/test/grafana"
)

func TestAddDataSource(t *testing.T) {
	ctx := context.Background()
	containerInfo, err := grafanacontainer.Start(ctx, t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	config := &Config{
		Host:          containerInfo.Host,
		Port:          strconv.Itoa(containerInfo.Port),
		AdminUser:     containerInfo.AdminUser,
		AdminPassword: containerInfo.AdminPassword,
	}
	g := New(config)

	err = g.AddDataSource("example", "mysql", srv.URL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestAddDashboard(t *testing.T) {
	ctx := context.Background()
	containerInfo, err := grafanacontainer.Start(ctx, t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	config := &Config{
		Host:          containerInfo.Host,
		Port:          strconv.Itoa(containerInfo.Port),
		AdminUser:     containerInfo.AdminUser,
		AdminPassword: containerInfo.AdminPassword,
	}
	g := New(config)

	dashboardJSON := []byte(`{"dashboard": {"title": "Sample Dashboard"}}`)
	err = g.AddDashboard(context.Background(), dashboardJSON)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
