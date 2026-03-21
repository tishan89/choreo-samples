/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-reading-list-web-app/internal/config"
	"io"
	"log"
	"net/http"
)

type Service struct {
	ApiUrl      string
	AccessToken string
}

// NewService creates a new instance of the service with the provided API URL
func NewService(apiUrl string, accessToken string) *Service {
	return &Service{
		ApiUrl:      apiUrl,
		AccessToken: accessToken,
	}
}

// FetchBooks fetches the list of books from the API using the configured access token
func (s *Service) FetchBooks() ([]config.Book, error) {
	client := &http.Client{}
	requestURL := s.ApiUrl + "/books"

	log.Printf("FetchBooks: preparing request url=%q token=%q", requestURL, maskToken(s.AccessToken))

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Printf("FetchBooks: request creation failed: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	log.Printf("FetchBooks: sending request method=%s url=%q auth_header_enabled=%t", req.Method, req.URL.String(), req.Header.Get("Authorization") != "")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("FetchBooks: request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FetchBooks: failed reading response body: %v", err)
		return nil, err
	}

	log.Printf("FetchBooks: response status=%s body=%q", resp.Status, truncateForLog(string(body), 512))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch books: %s", resp.Status)
	}

	var books []config.Book
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&books); err != nil {
		log.Printf("FetchBooks: failed decoding response body: %v", err)
		return nil, err
	}

	log.Printf("FetchBooks: decoded %d books", len(books))

	return books, nil
}

// AddNewBook adds a new book to the list using the configured access token
func (s Service) AddNewBook(book config.Book) error {
	client := &http.Client{}
	bookJson, err := json.Marshal(book)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.ApiUrl+"/books", bytes.NewBuffer(bookJson))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add book: %s", resp.Status)
	}
	return nil
}

// DeleteBook deletes a book from the list using the configured access token
func (s Service) DeleteBook(bookId string) error {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", s.ApiUrl+"/books/"+bookId, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete book: %s", resp.Status)
	}
	return nil
}

func maskToken(token string) string {
	if token == "" {
		return "<empty>"
	}
	if len(token) <= 8 {
		return token
	}

	return token[:4] + "..." + token[len(token)-4:]
}

func truncateForLog(value string, limit int) string {
	if len(value) <= limit {
		return value
	}

	return value[:limit] + "...(truncated)"
}
