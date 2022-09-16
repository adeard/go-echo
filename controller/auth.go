package controller

import (
	"go-echo/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
	jwtSecretKey           = "some-secret-key"
	jwtRefreshSecretKey    = "some-refresh-secret-key"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type TokenResult struct {
	RefreshToken string      `json:"refresh_token" form:"refresh_token"`
	AccessToken  string      `json:"access_token" form:"access_token"`
	Data         interface{} `json:"data"`
}

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

func GetJWTSecret() string {
	return jwtSecretKey
}

func GenerateTokensAndSetCookies(user *model.Users, c echo.Context) (interface{}, error) {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)
	// We generate here a new refresh token and saving it to the cookie.
	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	tokenResult := new(TokenResult)
	tokenResult.AccessToken = accessToken
	tokenResult.RefreshToken = refreshToken
	tokenResult.Data = user

	return tokenResult, nil
}

func generateRefreshToken(user *model.Users) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateAccessToken(user *model.Users) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(user *model.Users, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		Name: user.Nama,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// Purpose of this cookie is to store the user's name.
func setUserCookie(user *model.Users, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Nama
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(err error, c echo.Context) error {

	return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// If the user is not authenticated (no user token data in the context), don't do anything.
		if c.Get("user") == nil {
			return next(c)
		}
		// Gets user token from the context.
		u := c.Get("user").(*jwt.Token)

		claims := u.Claims.(*Claims)

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
			// Gets the refresh token from the cookie.
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				// Parses token and checks if it valid.
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					// If everything is good, update tokens.
					_, err = GenerateTokensAndSetCookies(&model.Users{Nama: claims.Name}, c)
				}
			}
		}

		return next(c)
	}
}

func SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Initiate a new User struct.
		u := new(model.Users)
		// Parse the submitted data and fill the User struct with the data from the SignIn form.
		if err := c.Bind(u); err != nil {
			// return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		fetchUser, error := model.GetOneByEmail(u.Email)
		if error != nil {
			return echo.NewHTTPError(http.StatusNotFound, error.Error())
		}

		// Compare the stored hashed password, with the hashed version of the password that was received.
		if err := bcrypt.CompareHashAndPassword([]byte(fetchUser.Password), []byte(u.Password)); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Password Incorrect")
		}
		// If password is correct, generate tokens and set cookies.
		result, err := GenerateTokensAndSetCookies(&fetchUser, c)

		if err != nil {

			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	}
}
