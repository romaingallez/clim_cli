package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Set_Clim(clim Clim) {
	// Manually build the query string to preserve order
	query := fmt.Sprintf("pow=%s&stemp=%s&mode=%s&shum=%s&f_rate=%s&f_dir=%s",
		url.QueryEscape(clim.Power),
		url.QueryEscape(clim.Temp),
		url.QueryEscape(clim.Mode),
		url.QueryEscape(clim.Shum),
		url.QueryEscape(clim.FanRate),
		url.QueryEscape(clim.FanDir),
	)

	// create an url.Url with clim.IP as host and "/set_control_info" as path
	u := &url.URL{
		Scheme:   "http",
		Host:     clim.IP,
		Path:     "aircon/set_control_info",
		RawQuery: query,
	}

	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyStr := string(body)
	fmt.Printf("Le serveur a répondu\nStatus: %s\nBody: %s\n", resp.Status, bodyStr)
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
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// GetControlInfo fetches control info from the given IP and returns it as a map.
func GetControlInfo(ip string) map[string]string {
	urlStr := fmt.Sprintf("http://%s/aircon/get_control_info", ip)
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("Error connecting to %s: %v\n", ip, err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response from %s: %v\n", ip, err)
		return nil
	}

	pairs := strings.Split(string(body), ",")
	parsedResponse := make(map[string]string)
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			parsedResponse[kv[0]] = kv[1]
		}
	}
	return parsedResponse
}

// FetchControlInfo is a context-aware variant returning an error on failure.
func FetchControlInfo(ctx context.Context, ip string) (map[string]string, error) {
	urlStr := fmt.Sprintf("http://%s/aircon/get_control_info", ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
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

// GetBasicInfo fetches basic info from the given IP and returns it as a map.
func GetBasicInfo(ip string) map[string]string {
	urlStr := fmt.Sprintf("http://%s/common/basic_info", ip)
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("Error connecting to %s: %v\n", ip, err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response from %s: %v\n", ip, err)
		return nil
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
	return parsedResponse
}

// unquote is a simple function to decode URL-encoded strings.
// In the Python code, it used `unquote`. Here, we're using Go's native function.
func unquote(s string) string {
	res, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return res
}
