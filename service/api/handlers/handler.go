package handlers

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	db "space-traders/repository/mysql"
	"space-traders/service/api/client"
	"space-traders/service/config"
)

type ViewHandler struct {
	userDB *db.Queries
	cfg    *config.Config
	Client *client.Client

	// paths that are not checked for a session cookie
	nonAuthPaths []string
	allowedPaths []string
}

func NewViewHandler(config *config.Config) *ViewHandler {
	c := &mysql.Config{
		User:   config.DBUser,
		Passwd: config.DBPass,
		Net:    config.DBNet,
		Addr:   config.DBAddr,
		DBName: config.DBName,
	}

	// create a connection pools
	pool, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	return &ViewHandler{
		userDB: db.New(pool),
		cfg:    config,
		Client: client.NewClient(config),
	}
}

// * should this be stored in a config file, so that it can be easily changed?
var nonAuthPaths = []string{
	"/login",
	"/login/submit",
	"/register",
	"/register/submit",
	"/com/header",
	"/com/footer",
}

// AddKeyToReq is a middleware that adds the API key to the request context.
// it ignores paths present in "nonAuthPaths".
func (vh *ViewHandler) AddKeyToReq() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, path := range nonAuthPaths {
				if c.Path() == path {
					return next(c)
				}
			}

			cookie, err := c.Cookie("session")
			if err != nil || cookie.Value == "" {
				return next(c)
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(401, "Invalid token")
				}

				return []byte(vh.cfg.JwtSecret), nil
			})
			if err != nil {
				next(c)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return next(c)
			}

			user, err := vh.userDB.GetUserWithAPIKeyByUsername(c.Request().Context(), sql.NullString{String: claims["username"].(string), Valid: true})
			if err != nil {
				return next(c)
			}

			c.Set("apiKey", user.ApiKey.String)
			c.Set("username", claims["username"])

			vh.Client.SetHeader("Authorization", "Bearer "+user.ApiKey.String)

			return next(c)
		}
	}
}

func (vh *ViewHandler) Favicon(c echo.Context) error {
	return c.File("service/static/favicon.ico")
}

// * Should this be stored in a config file, so that it can be easily changed?
var allowedPaths = []string{
	"/login",
	"/login/submit",
	"/register",
	"/register/submit",
}

// LoginRedirect is a middleware that checks if the user is logged in. If not, it redirects to the login page.
// It ignores paths present in "allowedPaths".
func (vh *ViewHandler) LoginRedirect() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, path := range allowedPaths {
				if c.Path() == path {
					return next(c)
				}
			}

			//should rename session to something like st_session to make it more unique
			cookie, err := c.Cookie("session")
			if err != nil || cookie.Value == "" {
				c.Redirect(302, "/login")
				return nil
			}

			return next(c)
		}
	}
}
