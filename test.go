package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const baseURL = "http://localhost:8080"

func TestRollDiceSingle(t *testing.T) {
	resp, err := http.Get(baseURL + "/rolldice/?player=TestPlayer")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Response: %s\n", string(body))
}

func TestRollDiceBestOfThree(t *testing.T) {
	resp, err := http.Get(baseURL + "/rolldice/?player=TestPlayer&mode=bestof3")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Response: %s\n", string(body))
}

func TestRollDiceBattle(t *testing.T) {
	resp, err := http.Get(baseURL + "/rolldice/?player=TestPlayer&mode=battle&opponent=OpponentPlayer")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Response: %s\n", string(body))
}

func TestViewLeaderboard(t *testing.T) {
	resp, err := http.Get(baseURL + "/leaderboard")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Leaderboard: %s\n", string(body))
}

func TestResetMetrics(t *testing.T) {
	resp, err := http.PostForm(baseURL+"/metrics/reset", url.Values{})
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Metrics Reset Response: %s\n", string(body))
}
