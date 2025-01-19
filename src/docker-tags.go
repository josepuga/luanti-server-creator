package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DockerHubResponse: Docker Hub API Response.
type DockerHubResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next string `json:"next"`
}

func getDockerTags(namespace string, image string) ([]string, error) {
	tags := []string{}
    baseURL := "https://hub.docker.com/v2/repositories/%s/%s/tags?page_size=100"

	url := fmt.Sprintf(baseURL, namespace, image)
	client := &http.Client{Timeout: 10 * time.Second}

	for url != "" {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching tags: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
		}

		var dockerResp DockerHubResponse
		err = json.NewDecoder(resp.Body).Decode(&dockerResp)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		for _, result := range dockerResp.Results {
			tags = append(tags, result.Name)
		}

		url = dockerResp.Next // Next page, if exists
	}
	return tags, nil
}

