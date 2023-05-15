package maven

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const TimeFormat = "20060102150405"

func FetchLatestVersion(client *http.Client, dependency string, repository string) (*Metadata, error) {
	split := strings.Split(dependency, ":")
	if len(split) < 2 {
		return nil, fmt.Errorf("invalid dependency: %s", dependency)
	}
	url := fmt.Sprintf("%s/%s/%s/maven-metadata.xml", repository, strings.ReplaceAll(split[0], ".", "/"), split[1])
	rs, err := client.Get(url)
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
	return &metadata, nil
}

type Metadata struct {
	GroupID    string     `xml:"groupId"`
	ArtifactID string     `xml:"artifactId"`
	Versioning Versioning `xml:"versioning"`
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
