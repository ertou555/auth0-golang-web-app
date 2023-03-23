package login

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/oauth2"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"01-Login/platform/authenticator"
)

// Handler for our login.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		invite := make([]oauth2.AuthCodeOption, 0)
		// 注册用户时必须要带audience，否则生成的token没有audience，即是Opaque access tokens无法使用
		audience := oauth2.SetAuthURLParam("audience", os.Getenv("AUTH0_AUDIENCE"))
		invite = append(invite, audience)

		// 被邀请人注册时要带organization和invitation发给auth0。
		organization := oauth2.SetAuthURLParam("organization", os.Getenv("AUTH0_CALLBACK_URL"))
		invitation := oauth2.SetAuthURLParam("invitation", "Qy3hMdYyFGEOdnERMNTb7GkBNtWWmwjm")

		invite = append(invite, organization)
		invite = append(invite, invitation)
		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state, invite...))
		//ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
