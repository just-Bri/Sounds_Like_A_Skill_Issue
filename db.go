// File: db.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

type PlayerScore struct {
	ID      *int   `json:"ID,omitempty"`
	Name    string `json:"Name"`
	Score   int    `json:"Score"`
	Version string `json:"Version"`
}

var SS string
var SU string

var Highscores []PlayerScore

func FetchHighscores() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", SU+"/rest/v1/highscores?select=*&order=Score.desc&limit=10", nil)
	if err != nil {
		return err
	}

	req.Header.Add("apikey", SS)
	req.Header.Add("Authorization", "Bearer "+SS)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &Highscores)
	if err != nil {
		return err
	}

	// Do I need this if I'm doing it in the query? /shrug
	sort.Slice(Highscores, func(i, j int) bool {
		return Highscores[i].Score > Highscores[j].Score
	})

	return nil
}

func LogHighscore() error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		playerScore := PlayerScore{
			Name:    PlayerInstance.Name,
			Score:   CurrentScore,
			Version: GameVersion,
		}

		jsonData, err := json.Marshal(playerScore)
		if err != nil {
			return fmt.Errorf("failed to marshal PlayerScore: %v", err)
		}

		req, err := http.NewRequest("POST", SU+"/rest/v1/highscores", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", SS)
		req.Header.Set("Authorization", "Bearer "+SS)
		req.Header.Set("Prefer", "return=minimal") // This tells Supabase to not return the inserted data

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			return fmt.Errorf("failed to send request: %v", err)
		} else {
			ScoreLogged = true
		}

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to insert highscore, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			return nil // Success
		}

		// If we get here, there was an error
		if i < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return fmt.Errorf("failed to insert highscore after %d attempts, last status: %d", maxRetries, resp.StatusCode)
	}

	return fmt.Errorf("failed to insert highscore after %d attempts", maxRetries)
}
