package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshToken string `json:"refresh_token"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExp        int64  `json:"at_exp"`
	RtExp        int64  `json:"rt_exp"`
}

func (n *TokenDetails) SendToken(c *gin.Context) {
	// Calculate max age in seconds
	accessMaxAge := int(time.Until(time.Unix(n.AtExp, 0)).Seconds())
	refreshMaxAge := int(time.Until(time.Unix(n.RtExp, 0)).Seconds())

	// Ensure max age is not negative
	if accessMaxAge < 0 {
		accessMaxAge = 0
	}
	if refreshMaxAge < 0 {
		refreshMaxAge = 0
	}

	// Set Session Token cookie using Gin's method
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"SessionToken", // name
		n.AccessToken,  // value
		accessMaxAge,   // maxAge in seconds
		"/",            // path
		// "",             // domain (empty means current domain)
		".songbooksofpraise.com", // domain (empty means current domain)
		true,                     // secure (HTTPS only)
		true,                     // httpOnly
	)

	// Set Refresh Token cookie
	c.SetCookie(
		"RefreshToken", // name
		n.RefreshToken, // value
		refreshMaxAge,  // maxAge in seconds
		"/",            // path
		// "",             // domain
		".songbooksofpraise.com", // domain
		true,                     // secure
		true,                     // httpOnly
	)
}

func UnsetToken(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)

	// Unset Session Token
	c.SetCookie(
		"SessionToken",
		"",
		-1, // negative maxAge deletes the cookie
		"/",
		"",
		true,
		true,
	)

	// Unset Refresh Token
	c.SetCookie(
		"RefreshToken",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
}

func CreateToken(user models.User) (TokenDetails, error) {
	tokenDetails := TokenDetails{}
	tokenDetails.AccessUUID = uuid.NewV4().String()
	tokenDetails.RefreshUUID = uuid.NewV4().String()

	// Creating Access Token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID
	tokenClaims["verified"] = user.Verified
	tokenClaims["access_uuid"] = tokenDetails.AccessUUID
	tokenClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.AtExp = tokenClaims["exp"].(int64)

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	access_token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return tokenDetails, err
	}
	tokenDetails.AccessToken = access_token

	// Creating Refresh Token
	tokenClaims = jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID
	tokenClaims["verified"] = user.Verified
	tokenClaims["refresh_uuid"] = tokenDetails.RefreshUUID
	tokenClaims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	tokenDetails.RtExp = tokenClaims["exp"].(int64)

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	refresh_token, err := rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return tokenDetails, err
	}
	tokenDetails.RefreshToken = refresh_token

	// Saving token in redis
	// client := redisdb.GetRedisConnection()
	// client.Set(tokenDetails.AccessUUID, user.ID, time.Until(time.Unix(tokenDetails.AtExp, 0)))
	// client.Set(tokenDetails.RefreshUUID, user.ID, time.Until(time.Unix(tokenDetails.RtExp, 0)))

	db := sqlite.GetDBConnection()
	db.Exec("INSERT OR REPLACE INTO session_tokens (access_uuid, refresh_uuid, user_id, at_exp, rt_exp) VALUES (?, ?, ?, ?, ?)",
		tokenDetails.AccessUUID, tokenDetails.RefreshUUID, user.ID, tokenDetails.AtExp, tokenDetails.RtExp)

	return tokenDetails, nil
}

func RetrieveUser(c *gin.Context) (models.User, error) {
	sessionTokenCookie, err := c.Request.Cookie("SessionToken")
	if err != nil {
		return models.User{}, err
	}

	// Verify Session Token is valid
	sessionToken, err := jwt.Parse(sessionTokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("session_token.invalid")
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return models.User{}, err
	}

	sessionTokenClaims := sessionToken.Claims.(jwt.MapClaims)
	user_id := int(sessionTokenClaims["user_id"].(float64))

	user, err := new(models.User).GetUserById(user_id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
