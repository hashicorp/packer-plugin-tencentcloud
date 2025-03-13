// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cvm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Profile struct {
	Type                string `json:"type,omitempty"`
	Region              string
	SecretId            string `json:"secretId,omitempty"`
	SecretKey           string `json:"secretKey,omitempty"`
	Token               string `json:"token,omitempty"`
	ExpiresAt           int64  `json:"expiresAt,omitempty"`
	RoleArn             string `json:"role-arn,omitempty"`
	RoleSessionName     string `json:"role-session-name,omitempty"`
	RoleSessionDuration int64  `json:"role-session-duration,omitempty"`
	Oauth               *Oauth `json:"oauth,omitempty"`
}

type Oauth struct {
	OpenId       string `json:"openId,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	ExpiresAt    int64  `json:"expiresAt,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Site         string `json:"site,omitempty"`
}

func getProfilePatch(cf *TencentCloudAccessConfig) (string, string, error) {
	var (
		profile              string
		sharedCredentialsDir string
		credentialPath       string
		configurePath        string
	)

	if cf.Profile != "" {
		profile = cf.Profile
	} else {
		profile = DEFAULT_PROFILE
	}

	if cf.SharedCredentialsDir != "" {
		sharedCredentialsDir = cf.SharedCredentialsDir
	}

	tmpSharedCredentialsDir, err := homedir.Expand(sharedCredentialsDir)
	if err != nil {
		return "", "", err
	}

	if tmpSharedCredentialsDir == "" {
		credentialPath = fmt.Sprintf("%s/.tccli/%s.credential", os.Getenv("HOME"), profile)
		configurePath = fmt.Sprintf("%s/.tccli/%s.configure", os.Getenv("HOME"), profile)
		if runtime.GOOS == "windows" {
			credentialPath = fmt.Sprintf("%s/.tccli/%s.credential", os.Getenv("USERPROFILE"), profile)
			configurePath = fmt.Sprintf("%s/.tccli/%s.configure", os.Getenv("USERPROFILE"), profile)
		}
	} else {
		credentialPath = fmt.Sprintf("%s/%s.credential", tmpSharedCredentialsDir, profile)
		configurePath = fmt.Sprintf("%s/%s.configure", tmpSharedCredentialsDir, profile)
	}

	return credentialPath, configurePath, nil
}

func loadConfigProfile(cf *TencentCloudAccessConfig) (*Profile, error) {
	var (
		credentialPath string
		configurePath  string
	)

	credentialPath, configurePath, err := getProfilePatch(cf)
	if err != nil {
		return nil, err
	}

	tcProfile := &Profile{}
	_, err = os.Stat(credentialPath)
	if !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(credentialPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, tcProfile)
		if err != nil {
			return nil, fmt.Errorf("credential file unmarshal failed, %s", err)
		}

		if tcProfile.Type == "oauth" {
			err := GetOauthConfig(tcProfile)
			if err != nil {
				return nil, fmt.Errorf("getOauthConfig failed, %v", err)
			}
		}
	} else {
		return nil, fmt.Errorf("please set a valid secret_id and secret_key or shared_credentials_dir, %s", err)
	}
	_, err = os.Stat(configurePath)
	if !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(configurePath)
		if err != nil {
			return nil, err
		}

		config := map[string]interface{}{}
		err = json.Unmarshal(data, &config)
		if err != nil {
			return nil, fmt.Errorf("configure file unmarshal failed, %s", err)
		}

	outerLoop:
		for k, v := range config {
			if k == "_sys_param" {
				tmpMap := v.(map[string]interface{})
				for tmpK, tmpV := range tmpMap {
					if tmpK == "region" {
						tcProfile.Region = strings.TrimSpace(tmpV.(string))
						break outerLoop
					}
				}
			}
		}
	} else {
		return nil, fmt.Errorf("please set a valid region or shared_credentials_dir, %s", err)
	}

	return tcProfile, nil
}
