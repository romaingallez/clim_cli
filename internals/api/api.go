package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyStr := string(body)
	// fmt.Println(resp.Status, resp.StatusCode)
	fmt.Printf("Le serveur a répondu\nStatus: %s\nBody: %s\n", resp.Status, bodyStr)
	// read the response body

}

// GetControlInfo fetches control info from the given IP and returns it as a map.
func GetControlInfo(ip string) map[string]string {
	url := fmt.Sprintf("http://%s/aircon/get_control_info", ip)
	resp, err := http.Get(url)
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

// GetBasicInfo fetches basic info from the given IP and returns it as a map.
func GetBasicInfo(ip string) map[string]string {
	url := fmt.Sprintf("http://%s/common/basic_info", ip)
	resp, err := http.Get(url)
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
		return s // return the original string if there's an error
	}
	return res
}
