package scoremaker

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type cookieRoundTripper struct {
	base   http.RoundTripper
	cookie string
}

func (c *cookieRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.cookie != "" {
		clone := req.Clone(req.Context())
		clone.Header.Set("Cookie", c.cookie)
		req = clone
	}
	return c.base.RoundTrip(req)
}

const (
	defaultUserAgent       = "UnityPlayer/2022.3.62f2"
	defaultPlatform        = "Android"
	defaultDeviceModel     = "samsung SM-G955F/dream2ltexx"
	defaultOperatingSystem = "Android"
	defaultUnityVersion    = "2022.3.62f2"
	defaultAcceptEncoding  = "deflate, gzip"
	defaultAccept          = "application/octet-stream"
)

var apiEndpoints = map[string]string{
	"jp": "https://production-game-api.sekai.colorfulpalette.org",
	"en": "https://n-production-game-api.sekai-en.com",
	"tw": "https://mk-zian-obt-cdn.bytedgame.com",
	"kr": "https://mkkorea-obt-prod01-cdn.bytedgame.com",
	"cn": "https://mkcn-prod-public-60001-1.dailygn.com",
}

type DownloadConfig struct {
	Region       string
	AppVersion   string
	AppHash      string
	IssueCookie  string
	SessionToken string
	Client       *http.Client
}

func BuildHeaders(cfg DownloadConfig) (http.Header, error) {
	if _, ok := apiEndpoints[strings.ToLower(cfg.Region)]; !ok {
		return nil, fmt.Errorf("unsupported region: %s", cfg.Region)
	}
	if strings.TrimSpace(cfg.AppVersion) == "" {
		return nil, fmt.Errorf("app version is required")
	}
	if strings.TrimSpace(cfg.AppHash) == "" {
		return nil, fmt.Errorf("app hash is required")
	}

	h := make(http.Header)
	h.Set("Accept", defaultAccept)
	h.Set("Content-Type", "application/octet-stream")
	h.Set("Accept-Encoding", defaultAcceptEncoding)
	h.Set("User-Agent", defaultUserAgent)
	h.Set("X-Platform", defaultPlatform)
	h.Set("X-DeviceModel", defaultDeviceModel)
	h.Set("X-OperatingSystem", defaultOperatingSystem)
	h.Set("X-Unity-Version", defaultUnityVersion)
	h.Set("X-App-Version", cfg.AppVersion)
	h.Set("X-App-Hash", cfg.AppHash)
	if strings.TrimSpace(cfg.SessionToken) != "" {
		h.Set("X-Session-Token", cfg.SessionToken)
	}
	if strings.TrimSpace(cfg.IssueCookie) != "" {
		h.Set("Cookie", cfg.IssueCookie)
	}

	return h, nil
}

func DownloadRawChart(scorePath string, cfg DownloadConfig) ([]byte, error) {
	region := strings.ToLower(cfg.Region)
	baseURL, ok := apiEndpoints[region]
	if !ok {
		return nil, fmt.Errorf("unsupported region: %s", cfg.Region)
	}
	if strings.TrimSpace(scorePath) == "" {
		return nil, fmt.Errorf("score path is required")
	}

	if strings.TrimSpace(cfg.IssueCookie) == "" {
		cookie, err := FetchIssueCookie(region, cfg.Client)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch issue cookie: %w", err)
		}
		cfg.IssueCookie = cookie
	}

	headers, err := BuildHeaders(cfg)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/blob/custom-music-score/full/%s", baseURL, strings.TrimPrefix(scorePath, "/"))

	baseTransport := http.DefaultTransport
	if cfg.Client != nil && cfg.Client.Transport != nil {
		baseTransport = cfg.Client.Transport
	}
	client := &http.Client{
		Transport: &cookieRoundTripper{
			base:   baseTransport,
			cookie: cfg.IssueCookie,
		},
	}
	if cfg.Client != nil && cfg.Client.Timeout != 0 {
		client.Timeout = cfg.Client.Timeout
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("download score failed: status=%d body=%q", resp.StatusCode, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
