package scoremaker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var issueEndpoints = map[string]string{
	"jp": "https://issue.sekai.colorfulpalette.org",
}

func normalizeIssueCookie(setCookieHeaders []string) string {
	if len(setCookieHeaders) == 0 {
		return ""
	}
	keys := []string{"CloudFront-Policy", "CloudFront-Signature", "CloudFront-Key-Pair-Id"}
	values := make(map[string]string, len(keys))

	for _, header := range setCookieHeaders {
		parts := strings.Split(header, ";")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			for _, key := range keys {
				prefix := key + "="
				if strings.HasPrefix(trimmed, prefix) {
					values[key] = trimmed
				}
			}
		}
	}

	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if v, ok := values[key]; ok {
			out = append(out, v)
		}
	}
	return strings.Join(out, "; ")
}

func FetchIssueCookie(region string, client *http.Client) (string, error) {
	issueBase, ok := issueEndpoints[strings.ToLower(strings.TrimSpace(region))]
	if !ok {
		return "", nil
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	req, err := http.NewRequest(http.MethodPost, issueBase+"/api/signature", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("issue endpoint returned status %d", resp.StatusCode)
	}
	cookie := normalizeIssueCookie(resp.Header.Values("Set-Cookie"))
	return cookie, nil
}

type AppInfo struct {
	AppVersion string
	AppHash    string
}

type appHashResponse struct {
	ProductionAndroid struct {
		AppHash    string `json:"app_hash"`
		AppVersion string `json:"app_version"`
	} `json:"production_android"`
	APKVersion string `json:"apk_version"`
}

func FetchAppInfo(region string, client *http.Client) (AppInfo, error) {
	region = strings.ToLower(strings.TrimSpace(region))
	if _, ok := apiEndpoints[region]; !ok {
		return AppInfo{}, fmt.Errorf("unsupported region: %s", region)
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/YangTheParrot/sekai-apphash/refs/heads/master/%s/apphash.json", region)
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return AppInfo{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return AppInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppInfo{}, fmt.Errorf("failed to fetch app info: status=%d", resp.StatusCode)
	}

	var data appHashResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return AppInfo{}, err
	}

	appVersion := strings.TrimSpace(data.ProductionAndroid.AppVersion)
	if region == "tw" || region == "kr" || region == "cn" {
		if v := strings.TrimSpace(data.APKVersion); v != "" {
			appVersion = v
		}
	}

	result := AppInfo{
		AppVersion: appVersion,
		AppHash:    strings.TrimSpace(data.ProductionAndroid.AppHash),
	}

	if result.AppVersion == "" || result.AppHash == "" {
		return AppInfo{}, fmt.Errorf("incomplete app info for region: %s", region)
	}

	return result, nil
}
