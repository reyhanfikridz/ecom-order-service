/*
Package middleware collection of middleware used for API
*/
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/reyhanfikridz/ecom-order-service/internal/config"
)

// User containing user data after authorization
type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

// AuthorizationMiddleware authorize each API route by checking JWT Token
func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get token
		token := GetTokenFromHeader(c.Request().Header)
		if token == "" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Token authorization empty/not found",
			})
		}

		// set form data
		formData := map[string]io.Reader{
			"token": strings.NewReader(token),
		}

		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		bFormDataWriter := multipart.NewWriter(&bFormData)
		for key, formDataReader := range formData {
			fieldWriter, err := bFormDataWriter.CreateFormField(key)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": err.Error(),
				})
			}

			_, err = io.Copy(fieldWriter, formDataReader)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": err.Error(),
				})
			}
		}
		bFormDataWriter.Close()

		// authorize to account service
		resp, err := http.Post(config.AccountServiceURL+"/api/authorize/",
			bFormDataWriter.FormDataContentType(),
			&bFormData)
		if err != nil { // if error occured
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}
		if resp.StatusCode != http.StatusOK { // if unauthorized
			return c.JSON(http.StatusForbidden, "Token authorization invalid")
		}

		// get user data from authorization response
		user, err := GetUserFromAuthorizationResp(resp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
		}

		c.Set("user", user)
		return next(c)
	}
}

// GetTokenFromHeader getting token (bearer) from request header
func GetTokenFromHeader(header http.Header) string {
	rawToken := header.Get("Authorization")
	if rawToken == "" {
		return ""
	}

	splitToken := strings.Split(rawToken, "Bearer ")
	if len(splitToken) <= 1 {
		return ""
	}

	token := splitToken[1]
	return token
}

// GetUserFromAuthorizationResp get user data from authorization response
func GetUserFromAuthorizationResp(resp *http.Response) (User, error) {
	user := User{}

	err := json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
