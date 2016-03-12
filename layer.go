package layer

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	base   = "https://api.layer.com"
	prefix = "apps"
)

var (
	// ErrMissingUserID __
	ErrMissingUserID = errors.New("Missing UserID")
	// ErrEmptyParticipants __
	ErrEmptyParticipants = errors.New("Empty Participants")
)

// Layer is an instance of a layer api object
type Layer struct {
	token   string
	appID   string
	version string
	timeout time.Duration
}

// Parameters contains the options passed in from the caller of request
type Parameters struct {
	Dedupe *string
	Path   string
	Body   []byte
}

// QueryParameters contains the possible query parameters to add onto a layer API call
type QueryParameters struct {
	PageSize int
	FromID   string
	SortBy   string
}

// NewLayer returns a new instance of a Layer struct
func NewLayer(token, appID, version string, timeout time.Duration) *Layer {
	return &Layer{
		token:   token,
		appID:   appID,
		version: version,
		timeout: timeout,
	}
}

func (l *Layer) request(method string, p *Parameters) (*http.Response, error) {
	method = strings.ToUpper(method)
	client := http.Client{Timeout: l.timeout}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s/%s/%s", base, prefix, l.appID, p.Path), bytes.NewBuffer(p.Body))
	if err != nil {
		return nil, err
	}

	if method == "PATCH" {
		req.Header.Set("Content-Type", "application/vnd.layer-patch+json")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	if p.Dedupe != nil {
		req.Header.Set("If-None-Match", *p.Dedupe)
	}

	req.Header.Set("Accept", fmt.Sprintf("application/vnd.layer+json; version=%s", l.version))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", l.token))

	return client.Do(req)

}
