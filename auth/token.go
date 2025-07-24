package auth

import (
	"fmt"
	"math"
	"os"
	"time"

	redisdb "github.com/BenjaminRA/himnario-backend/db/redis"
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
	c.Writer.Header().Add("Set-Cookie", fmt.Sprintf("SessionToken=%s; SameSite=none; Max-Age=%v; Path=/; Secure; HttpOnly", n.AccessToken, math.Floor(time.Until(time.Unix(n.AtExp, 0)).Seconds())))
	c.Writer.Header().Add("Set-Cookie", fmt.Sprintf("RefreshToken=%s; SameSite=none; Max-Age=%v; Path=/; Secure; HttpOnly", n.RefreshToken, math.Floor(time.Until(time.Unix(n.RtExp, 0)).Seconds())))
}

func UnsetToken(c *gin.Context) {
	c.Writer.Header().Add("Set-Cookie", "SessionToken=; SameSite=none; Max-Age=0; Path=/; Secure; HttpOnly")
	c.Writer.Header().Add("Set-Cookie", "RefreshToken=; SameSite=none; Max-Age=0; Path=/; Secure; HttpOnly")
}

func CreateToken(user models.User) (TokenDetails, error) {
	tokenDetails := TokenDetails{}
	tokenDetails.AccessUUID = uuid.NewV4().String()
	tokenDetails.RefreshUUID = uuid.NewV4().String()

	// Creating Access Token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID.Hex()
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
	tokenClaims["user_id"] = user.ID.Hex()
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
	client := redisdb.GetRedisConnection()
	client.Set(tokenDetails.AccessUUID, user.ID.Hex(), time.Until(time.Unix(tokenDetails.AtExp, 0)))
	client.Set(tokenDetails.RefreshUUID, user.ID.Hex(), time.Until(time.Unix(tokenDetails.RtExp, 0)))

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
	user_id := sessionTokenClaims["user_id"].(string)

	user, err := new(models.User).GetUserById(user_id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
