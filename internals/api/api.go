package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 4 * time.Second,
	Transport: &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         (&net.Dialer{Timeout: 3 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		TLSHandshakeTimeout: 3 * time.Second,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 10,
	},
}

// SetClim performs a control update with context and returns an error on failure.
func SetClim(ctx context.Context, clim Clim) error {
	query := fmt.Sprintf("pow=%s&stemp=%s&mode=%s&shum=%s&f_rate=%s&f_dir=%s",
		url.QueryEscape(clim.Power),
		url.QueryEscape(clim.Temp),
		url.QueryEscape(clim.Mode),
		url.QueryEscape(clim.Shum),
		url.QueryEscape(clim.FanRate),
		url.QueryEscape(clim.FanDir),
	)
	u := &url.URL{Scheme: "http", Host: clim.IP, Path: "aircon/set_control_info", RawQuery: query}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil http response")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// FetchControlInfo fetches control info using context and returns a parsed map.
func FetchControlInfo(ctx context.Context, ip string) (map[string]string, error) {
	urlStr := fmt.Sprintf("http://%s/aircon/get_control_info", ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("nil http response")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pairs := strings.Split(string(body), ",")
	parsedResponse := make(map[string]string)
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			parsedResponse[kv[0]] = kv[1]
		}
	}
	return parsedResponse, nil
}

// FetchBasicInfo fetches basic info using context and returns a parsed map.
func FetchBasicInfo(ctx context.Context, ip string) (map[string]string, error) {
	urlStr := fmt.Sprintf("http://%s/common/basic_info", ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("nil http response")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pairs := strings.Split(string(body), ",")
	parsedResponse := make(map[string]string)
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key, value := kv[0], kv[1]
			if key == "name" || key == "grp_name" {
				parsedResponse[key] = unquote(value)
			} else {
				parsedResponse[key] = value
			}
		}
	}
	return parsedResponse, nil
}

// Legacy helpers (deprecated): retained for backward-compatibility
func GetControlInfo(ip string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	m, _ := FetchControlInfo(ctx, ip)
	return m
}

func GetBasicInfo(ip string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	m, _ := FetchBasicInfo(ctx, ip)
	return m
}

// unquote is a simple function to decode URL-encoded strings.
func unquote(s string) string {
	res, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return res
}
