package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
)

type logEntry struct {
	Entry       *log.Entry
	Timestamp   time.Time
	ServiceName string
	ModuleName  string
}

// LokiHook is a custom Logrus hook for sending logs to Loki and interacting with its API.
type LokiHook struct {
	url          string
	labels       map[string]string
	entries      chan logEntry
	sendInterval time.Duration
	maxBatchSize int
	retryCount   int
	retryDelay   time.Duration
}

// NewLokiHook initializes a new LokiHook with the specified configuration.
func NewLokiHook(cfg *Config) *LokiHook {
	hook := &LokiHook{
		url:          cfg.LokiURL,
		labels:       cfg.Labels,
		entries:      make(chan logEntry, cfg.MaxBufferSize),
		sendInterval: cfg.SendInterval,
		maxBatchSize: cfg.MaxBatchSize,
		retryCount:   cfg.RetryCount,
		retryDelay:   cfg.RetryDelay,
	}
	go hook.batchSendLogs()
	return hook
}

// pushAction returns the URL for pushing logs to Loki.
func (hook *LokiHook) pushAction() string {
	return hook.url + "/loki/api/v1/push"
}

// queryRangeAction returns the URL for querying logs from Loki.
func (hook *LokiHook) queryRangeAction(query string, limit int, start, end time.Time) string {
	return fmt.Sprintf("%s/loki/api/v1/query_range?query=%s&limit=%d&start=%d&end=%d",
		hook.url, query, limit, start.UnixNano(), end.UnixNano())
}

// labelsAction returns the URL for retrieving labels from Loki.
func (hook *LokiHook) labelsAction() string {
	return hook.url + "/loki/api/v1/labels"
}

// Fire queues a log entry to be sent to Loki.
func (hook *LokiHook) Fire(entry *log.Entry) error {
	serviceName, _ := entry.Data["serviceName"].(string)
	moduleName, _ := entry.Data["moduleName"].(string)

	hook.entries <- logEntry{
		Entry:       entry,
		Timestamp:   time.Now(),
		ServiceName: serviceName,
		ModuleName:  moduleName,
	}
	return nil
}

// Levels returns the log levels for which this hook is activated.
func (hook *LokiHook) Levels() []log.Level {
	return log.AllLevels
}

// batchSendLogs processes and sends batches of log entries to Loki.
func (hook *LokiHook) batchSendLogs() {
	ticker := time.NewTicker(hook.sendInterval)
	defer ticker.Stop()

	var logs []logEntry
	for {
		select {
		case entry := <-hook.entries:
			logs = append(logs, entry)
			if len(logs) >= hook.maxBatchSize {
				hook.sendBatchLogs(logs)
				logs = logs[:0]
			}
		case <-ticker.C:
			if len(logs) > 0 {
				hook.sendBatchLogs(logs)
				logs = logs[:0]
			}
		}
	}
}

// sendBatchLogs prepares and sends log entries to Loki.
func (hook *LokiHook) sendBatchLogs(logs []logEntry) {
	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": hook.labels,
				"values": formatLogEntries(logs),
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling logs: %v\n", err)
		return
	}

	operation := func() error {
		req, err := http.NewRequest("POST", hook.pushAction(), bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return errors.New("failed to send logs to Loki")
		}
		return nil
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = hook.retryDelay
	backoffConfig.MaxElapsedTime = time.Duration(hook.retryCount) * hook.retryDelay

	if err := backoff.Retry(operation, backoffConfig); err != nil {
		fmt.Printf("Failed to send logs to Loki after multiple attempts: %v\n", err)
	}
}

// formatLogEntries converts log entries into a format suitable for Loki.
func formatLogEntries(entries []logEntry) [][]string {
	formatted := make([][]string, len(entries))
	for i, entry := range entries {
		formatted[i] = []string{
			fmt.Sprintf("%d", entry.Timestamp.UnixNano()),
			formatLogEntry(entry),
		}
	}
	return formatted
}

// formatLogEntry formats a single log entry into a JSON string.
func formatLogEntry(entry logEntry) string {
	attributes := map[string]interface{}{
		"level":       entry.Entry.Level.String(),
		"time":        entry.Timestamp.Format(time.RFC3339Nano),
		"message":     entry.Entry.Message,
		"serviceName": entry.ServiceName,
		"moduleName":  entry.ModuleName,
	}

	for k, v := range entry.Entry.Data {
		attributes[k] = v
	}

	entryJSON, err := json.Marshal(attributes)
	if err != nil {
		return ""
	}
	return string(entryJSON)
}

// QueryRange executes a basic query against the Loki instance.
func (hook *LokiHook) QueryRange(query string, limit int, start, end time.Time) ([]byte, error) {
	resp, err := http.Get(hook.queryRangeAction(query, limit, start, end))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to execute query: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// GetLabels fetches available labels from Loki.
func (hook *LokiHook) GetLabels() ([]byte, error) {
	resp, err := http.Get(hook.labelsAction())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get labels: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
