package api

import (
	"fmt"
	"net/http"
	"net/url"
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

	resp, err := http.Get(u.String())

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(resp.Status, resp.StatusCode)
	fmt.Printf("Le serveur a répondu\nStatus: %s\n", resp.Status)

}
