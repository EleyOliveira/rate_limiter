package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestConcurrentRequests(t *testing.T) {

	numRequests := 10

	for i := 0; i < numRequests; i++ {
		req, err := http.NewRequest("GET", "http://localhost:8080", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("API_KEY", fmt.Sprintf(`my-api-key=%d`, i+1))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}

	}
}
