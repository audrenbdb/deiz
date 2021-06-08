package contact

import (
	"github.com/audrenbdb/deiz/email"
	"github.com/audrenbdb/deiz/psql"
	"github.com/labstack/echo"
	"net/http"
)

const deizContactAddress = "contact@deiz.fr"

type contactForm struct {
	ClinicianID int    `json:"clinicianId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Message     string `json:"message"`
}

type getInTouchForm struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Job   string `json:"job"`
	City  string `json:"city"`
}

func echoHandlePostContactForm(send sendContactForm) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var f contactForm
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := send(ctx, f); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func echoHandlePostGetInTouchForm(send sendGetInTouchForm) echo.HandlerFunc {
	return func(c echo.Context) error {
		var f getInTouchForm
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := send(f); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func RegisterService(e *echo.Echo, db psql.PGX, send email.Send) {
	repo := psqlRepo{db: db}
	sendContactForm := sendContactFormFn(repo.createGetClinicianByIDFunc(), send)
	sendGetInTouchForm := sendGetInTouchFormFn(send)

	e.POST("/api/forms/contact-forms", echoHandlePostContactForm(sendContactForm))
	e.POST("/api/forms/get-in-touch-forms", echoHandlePostGetInTouchForm(sendGetInTouchForm))
}
