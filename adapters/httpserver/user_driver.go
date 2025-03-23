package httpserver

import (
	"bytes"
	"encoding/json"
	"mytodoapp/adapters/httpserver/handler"
	"net/http"
)

type UserDriver struct {
	BaseURL string
	Client  *http.Client
	Token   string
}

func (d *UserDriver) RegisterUser(email string, password string) error {
	buffer := new(bytes.Buffer)
	payload := handler.RegisterUserPayload{Email: email, Password: password}
	json.NewEncoder(buffer).Encode(&payload)
	res, err := d.Client.Post(d.BaseURL+"/register", "application/json", buffer)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

func (d *UserDriver) LoginUser(email string, password string) (token string, error error) {
	buffer := new(bytes.Buffer)
	payload := handler.LoginUserPayload{Email: email, Password: password}
	json.NewEncoder(buffer).Encode(&payload)
	res, err := d.Client.Post(d.BaseURL+"/login", "application/json", buffer)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var result LoginUserResponse
	json.NewDecoder(res.Body).Decode(&result)
	// set cookie on client
	d.Token = result.Token

	return result.Token, nil
}

func (d *UserDriver) GetUserProfile(token string) (handler.UserProfilePayload, error) {
	req, err := http.NewRequest(http.MethodGet, d.BaseURL+"/profile", nil)
	if err != nil {
		return handler.UserProfilePayload{}, err
	}
	req.Header.Add("Authorization", d.Token)
	res, err := d.Client.Do(req)
	if err != nil {
		return handler.UserProfilePayload{}, err
	}
	defer res.Body.Close()

	var result handler.UserProfilePayload
	json.NewDecoder(res.Body).Decode(&result)

	return result, nil
}
