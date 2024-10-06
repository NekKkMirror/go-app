package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Config struct {
	Host             string `json:"host"`
	Port             string `json:"port"`
	AdminUser        string `json:"adminUser"`
	AdminPassword    string `json:"adminPassword"`
	ProvisioningPath string `json:"provisioningPath"`
	DashboardsPath   string `json:"dashboardsPath"`
}

type Grafana struct {
	Host          string
	Port          string
	AdminUser     string
	AdminPassword string
}

type Folder struct {
	Id    int64  `json:"id"`
	Uid   string `json:"uid"`
	Title string `json:"title"`
}

func New(config *Config) *Grafana {
	return &Grafana{
		Host:          config.Host,
		Port:          config.Port,
		AdminUser:     config.AdminUser,
		AdminPassword: config.AdminPassword,
	}
}

func (g *Grafana) URL() string {
	return fmt.Sprintf("http://%s:%s", g.Host, g.Port)
}

// AddDataSource adds a new data source to Grafana.
func (g *Grafana) AddDataSource(name, sourceType, url string) error {
	dataSourcePayload := fmt.Sprintf(`{
		"name": "%s",
		"type": "%s",
		"url": "%s",
		"access": "proxy",
		"basicAuth": false
	}`, name, sourceType, url)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/datasources", g.URL()), bytes.NewBuffer([]byte(dataSourcePayload)))
	if err != nil {
		return errors.Wrap(err, "failed to create request for Grafana")
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to add data source to Grafana")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add data source in Grafana, status: %s", resp.Status)
	}

	return nil
}

// UpdateDataSource updates an existing data source in Grafana.
func (g *Grafana) UpdateDataSource(id int64, dataSourcePayload string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/datasources/%d", g.URL(), id), bytes.NewBuffer([]byte(dataSourcePayload)))
	if err != nil {
		return errors.Wrap(err, "failed to create request for Grafana")
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to update data source to Grafana")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update data source in Grafana, status: %s", resp.Status)
	}

	return nil
}

// DeleteDataSource removes a data source from Grafana by ID.
func (g *Grafana) DeleteDataSource(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/datasources/%d", g.URL(), id), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create delete request for Grafana")
	}
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete data source from Grafana")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete data source in Grafana, status: %s", resp.Status)
	}

	return nil
}

// AddDashboard uploads a new dashboard to Grafana.
func (g *Grafana) AddDashboard(ctx context.Context, dashboardJSON []byte) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/dashboards/db", g.URL()), bytes.NewBuffer(dashboardJSON))
	if err != nil {
		return errors.Wrap(err, "failed to create request for Grafana")
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to add dashboard to Grafana")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add dashboard in Grafana, status: %s", resp.Status)
	}

	return nil
}

// ExportDashboard exports a dashboard from Grafana.
func (g *Grafana) ExportDashboard(uid string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/dashboards/uid/%s", g.URL(), uid), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request for exporting dashboard")
	}
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to export dashboard from Grafana")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to export dashboard in Grafana, status: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

// CreateFolder creates a new folder for dashboards in Grafana.
func (g *Grafana) CreateFolder(ctx context.Context, title string) (*Folder, error) {
	folderPayload := map[string]string{"title": title}
	payloadBytes, err := json.Marshal(folderPayload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal folder payload")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/folders", g.URL()), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request for creating folder in Grafana")
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create folder in Grafana")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create folder in Grafana, status: %s", resp.Status)
	}

	var folder Folder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return nil, errors.Wrap(err, "failed to decode folder response")
	}

	return &folder, nil
}

// MonitorHealth checks the health of the Grafana instance.
func (g *Grafana) MonitorHealth() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/health", g.URL()), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request for health check")
	}
	req.SetBasicAuth(g.AdminUser, g.AdminPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to check Grafana health")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Grafana is not healthy, status: %s", resp.Status)
	}

	return nil
}

// TODO TESTS
