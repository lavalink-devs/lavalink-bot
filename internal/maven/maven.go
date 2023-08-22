package maven

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	TimeFormat = "20060102150405"
	CacheTTL   = time.Minute * 5
)

func New(client *http.Client) *Maven {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Maven{
		client: client,
		cancel: cancel,
		cache:  make(map[string]*Metadata),
	}

	go m.cleanupCache(ctx)

	return m
}

type Maven struct {
	client *http.Client
	cancel context.CancelFunc

	cacheMu sync.Mutex
	cache   map[string]*Metadata
}

func (m *Maven) Close() {
	m.client.CloseIdleConnections()
	m.cancel()
}

func (m *Maven) cleanupCache(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			m.cacheMu.Lock()
			for k, v := range m.cache {
				if time.Since(v.FetchedAt) > CacheTTL {
					delete(m.cache, k)
				}
			}
			m.cacheMu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (m *Maven) FetchLatestVersion(dependency string, repository string) (*Metadata, error) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()
	if metadata, ok := m.cache[dependency]; ok {
		return metadata, nil
	}

	split := strings.Split(dependency, ":")
	if len(split) < 2 {
		return nil, fmt.Errorf("invalid dependency: %s", dependency)
	}
	url := fmt.Sprintf("%s/%s/%s/maven-metadata.xml", repository, strings.ReplaceAll(split[0], ".", "/"), split[1])
	rs, err := m.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get maven-metadata.xml: %w", err)
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get maven-metadata.xml: %s", rs.Status)
	}

	var metadata Metadata
	decoder := xml.NewDecoder(rs.Body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	if err = decoder.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode maven-metadata.xml: %w", err)
	}
	metadata.FetchedAt = time.Now()
	m.cache[dependency] = &metadata
	return &metadata, nil
}

type Metadata struct {
	GroupID    string     `xml:"groupId"`
	ArtifactID string     `xml:"artifactId"`
	Versioning Versioning `xml:"versioning"`
	FetchedAt  time.Time  `xml:"-"`
}

func (m Metadata) Latest() string {
	if m.Versioning.Latest != "" {
		return m.Versioning.Latest
	}
	if m.Versioning.Release != "" {
		return m.Versioning.Release
	}
	if len(m.Versioning.Versions) == 0 {
		return "unknown"
	}
	return m.Versioning.Versions[len(m.Versioning.Versions)-1]
}

type Versioning struct {
	Latest      string   `xml:"latest"`
	Release     string   `xml:"release"`
	Versions    []string `xml:"versions>version"`
	LastUpdated Time     `xml:"lastUpdated"`
}

type Time string

func (t Time) Time() (time.Time, error) {
	return time.Parse(TimeFormat, string(t))
}
