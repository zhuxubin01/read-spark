package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DictionaryService struct {
	client  *http.Client
	baseURL string
}

type DictionaryResult struct {
	Word     string   `json:"word"`
	Phonetic string   `json:"phonetic,omitempty"`
	Meanings []string `json:"meanings"`
}

func NewDictionaryService() *DictionaryService {
	return NewDictionaryServiceWithClient("https://api.dictionaryapi.dev", &http.Client{Timeout: 8 * time.Second})
}

func NewDictionaryServiceWithClient(baseURL string, client *http.Client) *DictionaryService {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.dictionaryapi.dev"
	}
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	return &DictionaryService{client: client, baseURL: strings.TrimRight(baseURL, "/")}
}

func (s *DictionaryService) Lookup(ctx context.Context, word string) (*DictionaryResult, error) {
	w := strings.TrimSpace(strings.ToLower(word))
	if w == "" {
		return nil, fmt.Errorf("word is required")
	}

	endpoint := s.baseURL + "/api/v2/entries/en/" + url.PathEscape(w)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("word not found")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dictionary upstream error: %d", resp.StatusCode)
	}

	var entries []struct {
		Word      string `json:"word"`
		Phonetic  string `json:"phonetic"`
		Phonetics []struct {
			Text string `json:"text"`
		} `json:"phonetics"`
		Meanings []struct {
			Definitions []struct {
				Definition string `json:"definition"`
			} `json:"definitions"`
		} `json:"meanings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("word not found")
	}

	result := &DictionaryResult{Word: entries[0].Word}
	if entries[0].Phonetic != "" {
		result.Phonetic = entries[0].Phonetic
	} else {
		for _, p := range entries[0].Phonetics {
			if strings.TrimSpace(p.Text) != "" {
				result.Phonetic = p.Text
				break
			}
		}
	}

	for _, m := range entries[0].Meanings {
		for _, d := range m.Definitions {
			if strings.TrimSpace(d.Definition) == "" {
				continue
			}
			result.Meanings = append(result.Meanings, d.Definition)
			if len(result.Meanings) >= 6 {
				return result, nil
			}
		}
	}

	return result, nil
}
