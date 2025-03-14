// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const _API_ENDPOINT = "https://cli.cloud.tencent.com"

func GetOauthConfig(p *Profile) error {
	if p.Oauth == nil || p.Oauth.RefreshToken == "" || p.Oauth.OpenId == "" {
		return fmt.Errorf("Oauth authentication information is not configured correctly")
	}
	client := NewAPIClient()

	expired := false
	if p.Oauth.ExpiresAt != 0 {
		now := time.Now()
		futureTime := now.Add(30 * time.Second)
		targetTime := time.Unix(p.Oauth.ExpiresAt, 0)
		if futureTime.After(targetTime) {
			expired = true
		}
	}

	if expired {
		response, err := client.RefreshUserToken(p.Oauth.RefreshToken, p.Oauth.OpenId, p.Oauth.Site)
		if err != nil {
			return err
		}
		if response != nil {
			if response.AccessToken != "" {
				p.Oauth.AccessToken = response.AccessToken
			}
			if response.ExpiresAt != 0 {
				p.Oauth.ExpiresAt = response.ExpiresAt
			}
		}
	}

	// 获取临时token
	response, err := client.GetThirdPartyFederationToken(p.Oauth.AccessToken, p.Oauth.Site)
	if err != nil {
		return err
	}
	if response != nil {
		if response.SecretId != "" {
			p.SecretId = response.SecretId
		}
		if response.SecretKey != "" {
			p.SecretKey = response.SecretKey
		}
		if response.Token != "" {
			p.Token = response.Token
		}
	}

	return nil
}

type APIClient struct {
	Client *http.Client
}

// 创建新的APIClient
func NewAPIClient() *APIClient {
	return &APIClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetThirdPartyFederationToken Obtaining a temporary user certificate
func (c *APIClient) GetThirdPartyFederationToken(accessToken, site string) (*GetTempCredResponse, error) {
	apiEndpoint := _API_ENDPOINT + "/get_temp_cred"
	traceId := uuid.New().String()

	body := GetTempCredRequest{
		TraceId:     traceId,
		AccessToken: accessToken,
		Site:        site,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	response := &GetTempCredResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("get_temp_cred: %s", response.Error)
	}

	return response, err
}

// RefreshUserToken Refresh user third-party access_token
func (c *APIClient) RefreshUserToken(refToken, openId, site string) (*RefreshTokenResponse, error) {
	apiEndpoint := _API_ENDPOINT + "/refresh_user_token"

	traceId := uuid.New().String()

	body := RefreshTokenRequest{
		TraceId:      traceId,
		RefreshToken: refToken,
		OpenId:       openId,
		Site:         site,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	response := &RefreshTokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("refresh_user_token: %s", response.Error)
	}
	return response, nil
}

type RefreshTokenRequest struct {
	TraceId      string `json:"TraceId"`
	RefreshToken string `json:"RefreshToken"`
	OpenId       string `json:"OpenId"`
	Site         string `json:"Site"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"AccessToken"`
	ExpiresAt   int64  `json:"ExpiresAt"`
	Error       string `json:"Error"`
}

type GetTempCredRequest struct {
	TraceId     string `json:"TraceId"`
	AccessToken string `json:"AccessToken"`
	Site        string `json:"Site"`
}

type GetTempCredResponse struct {
	SecretId  string `json:"SecretId"`
	SecretKey string `json:"SecretKey"`
	Token     string `json:"Token"`
	ExpiresAt int64  `json:"ExpiresAt"`
	Error     string `json:"Error"`
}
