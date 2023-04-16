package wikipedia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WikiResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID  int64  `json:"pageid"`
			NS      int64  `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func fetchWikipediaSummary(topic string) (string, error) {
	apiUrl := "https://en.wikipedia.org/w/api.php"

	queryParams := url.Values{}
	queryParams.Set("action", "query")
	queryParams.Set("format", "json")
	queryParams.Set("prop", "extracts")
	queryParams.Set("exintro", "")
	queryParams.Set("explaintext", "")
	queryParams.Set("titles", topic)

	resp, err := http.Get(fmt.Sprintf("%s?%s", apiUrl, queryParams.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to fetch Wikipedia data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var wikiResponse WikiResponse
	err = json.Unmarshal(body, &wikiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse Wikipedia JSON: %v", err)
	}

	for _, page := range wikiResponse.Query.Pages {
		return page.Extract, nil
	}

	return "", fmt.Errorf("no information found for topic: %s", topic)
}
