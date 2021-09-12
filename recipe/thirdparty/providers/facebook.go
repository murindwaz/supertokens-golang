package providers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

const facebookID = "facebook"

func Facebook(config FacebookConfig) tpmodels.TypeProvider {
	return tpmodels.TypeProvider{
		ID: facebookID,
		Get: func(redirectURI, authCodeFromRequest *string) tpmodels.TypeProviderGetResponse {
			accessTokenAPIURL := "https://graph.facebook.com/v9.0/oauth/access_token"
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": config.ClientSecret,
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://www.facebook.com/v9.0/dialog/oauth"
			scopes := []string{"email"}
			if config.Scope != nil {
				scopes = append(scopes, config.Scope...)
			}

			authorizationRedirectParams := map[string]interface{}{
				"scope":         strings.Join(scopes, " "),
				"response_type": "code",
				"client_id":     config.ClientID,
			}

			return tpmodels.TypeProviderGetResponse{
				AccessTokenAPI: tpmodels.AccessTokenAPI{
					URL:    accessTokenAPIURL,
					Params: accessTokenAPIParams,
				},
				AuthorisationRedirect: tpmodels.AuthorisationRedirect{
					URL:    authorisationRedirectURL,
					Params: authorizationRedirectParams,
				},
				GetProfileInfo: func(authCodeResponse interface{}) (tpmodels.UserInfo, error) {
					authCodeResponseJson, err := json.Marshal(authCodeResponse)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					var accessTokenAPIResponse facebookGetProfileInfoInput
					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					accessToken := accessTokenAPIResponse.AccessToken
					response, err := getFacebookAuthRequest(accessToken)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}
					userInfo := response.(map[string]interface{})
					ID := userInfo["id"].(string)
					email := userInfo["email"].(string)
					if email == "" {
						return tpmodels.UserInfo{
							ID: ID,
						}, nil
					}
					isVerified := userInfo["verified_email"].(bool)
					return tpmodels.UserInfo{
						ID: ID,
						Email: &tpmodels.EmailStruct{
							ID:         email,
							IsVerified: isVerified,
						},
					}, nil
				},
			}
		},
	}
}

func getFacebookAuthRequest(accessToken string) (interface{}, error) {
	url := "https://graph.facebook.com/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("access_token", accessToken)
	req.Header.Add("fields", "id,email")
	req.Header.Add("format", "json")
	return doGetRequest(req)
}

type facebookGetProfileInfoInput struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
