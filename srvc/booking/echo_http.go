package booking

import (
	"github.com/audrenbdb/deiz/auth"
	"github.com/labstack/echo/v4"
)

type echoServer struct {
	router *echo.Echo
	auth   auth.CredentialsGetter
	repo   *repo
}

func (s *echoServer) HandleGetBookings(c echo.Context) error {
	return nil
}

func (s *echoServer) RegisterService() {

}
