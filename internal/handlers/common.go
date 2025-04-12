package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResp struct {
	Message string `json:"message"`
}

type OAuthInfo struct {
	Username string
	Email    string
}

type OAuthInfoExtractor func(token string) (string, string, error)

func GithubInfoExtractor(token string) (string, string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error getting user info, estado: %s", resp.Status)
	}

	type GithubUser struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}

	var user GithubUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", "", err
	}

	if user.Email == "" {
		req2, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		if err != nil {
			return "", "", err
		}
		req2.Header.Set("Authorization", "Bearer "+token)
		req2.Header.Set("Accept", "application/json")

		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			return "", "", err
		}
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			return "", "", fmt.Errorf("error getting emails, status: %s", resp2.Status)
		}

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		if err = json.NewDecoder(resp2.Body).Decode(&emails); err != nil {
			return "", "", err
		}

		for _, e := range emails {
			if e.Primary {
				user.Email = e.Email
				break
			}
		}

		if user.Email == "" && len(emails) > 0 {
			user.Email = emails[0].Email
		}
	}

	return user.Login, user.Email, nil
}

func GoogleInfoExtractor(token string) (string, string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error gtting user info, status: %s", resp.Status)
	}

	type GoogleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
	}

	var user GoogleUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", "", err
	}

	return user.Name, user.Email, nil
}
