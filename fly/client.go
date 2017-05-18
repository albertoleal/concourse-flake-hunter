package fly

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/event"
	"github.com/concourse/go-concourse/concourse"
	"golang.org/x/oauth2"
)

//go:generate counterfeiter . Client
type Client interface {
	ConcourseURL() string
	Builds(concourse.Page) ([]atc.Build, concourse.Pagination, error)
	BuildEvents(buildID string) ([]byte, error)
}

type client struct {
	concourseURL string
	username     string
	password     string
	team         string

	concourseCli concourse.Client
}

func NewClient(concourseURL, username, password, team string) *client {
	c := &client{
		concourseURL: concourseURL,
		username:     username,
		password:     password,
		team:         team,
	}
	return c
}

func (c *client) ConcourseURL() string {
	return c.concourseURL
}

func (c *client) Builds(page concourse.Page) ([]atc.Build, concourse.Pagination, error) {
	client, err := c.concourseClient()
	if err != nil {
		return []atc.Build{}, concourse.Pagination{}, err
	}
	return client.Builds(page)
}

func (c *client) BuildEvents(buildID string) ([]byte, error) {
	client, err := c.concourseClient()
	if err != nil {
		return []byte{}, err
	}
	events, err := client.BuildEvents(buildID)
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer([]byte{})
	var buildConfig event.TaskConfig
	for {
		ev, err := events.NextEvent()
		if err != nil {
			if err == io.EOF {
				return buf.Bytes(), nil
			} else {
				panic("failed to parse event")
			}
		}

		switch e := ev.(type) {
		case event.Log:
			fmt.Fprintf(buf, "%s", e.Payload)

		case event.InitializeTask:
			buildConfig = e.TaskConfig

		case event.StartTask:
			argv := strings.Join(append([]string{buildConfig.Run.Path}, buildConfig.Run.Args...), " ")
			fmt.Fprintf(buf, "%s\n", argv)

		case event.Error:
			fmt.Fprintf(buf, "%s\n", e.Message)
		}
	}
	return buf.Bytes(), events.Close()
}

func (c *client) concourseClient() (concourse.Client, error) {
	if c.concourseCli != nil {
		return c.concourseCli, nil
	}

	httpClient := &http.Client{
		Transport: basicAuthTransport{
			username: c.username,
			password: c.password,
			base:     transport(),
		},
	}

	client := concourse.NewClient(c.concourseURL, httpClient)
	t := client.Team(c.team)
	token, err := t.AuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate to team: %s", err)
	}

	oAuthToken := &oauth2.Token{
		TokenType:   token.Type,
		AccessToken: token.Value,
	}

	transport := transport()
	transport = &oauth2.Transport{
		Source: oauth2.StaticTokenSource(oAuthToken),
		Base:   transport,
	}

	c.concourseCli = concourse.NewClient(c.concourseURL, &http.Client{Transport: transport})
	return c.concourseCli, nil
}

type basicAuthTransport struct {
	username string
	password string

	base http.RoundTripper
}

func (t basicAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.username, t.password)
	return t.base.RoundTrip(r)
}

func transport() http.RoundTripper {
	var transport http.RoundTripper

	transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		Proxy: http.ProxyFromEnvironment,
	}

	return transport
}
