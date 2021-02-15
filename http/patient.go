package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

type (
	PatientSearcher interface {
		SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error)
	}
	PatientAdder interface {
		AddPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
)

func handleGetPatients(searcher PatientSearcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		search := c.QueryParam("search")
		if search == "" {
			return c.JSON(http.StatusBadRequest, "no search provided")
		}
		p, err := searcher.SearchPatient(ctx, search, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, p)
	}
}

func handlePostPatient(adder PatientAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var p deiz.Patient
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, errBind)
		}
		if err := adder.AddPatient(ctx, &p, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, p)
	}
}

/*
func handleGetPatients(searchPatients deiz.SearchPatients) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		search := c.QueryParam("search")
		if search == "" {
			return c.JSON(http.StatusBadRequest, "no search provided")
		}
		p, err := searchPatients(ctx, search, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, p)
	}
}

func handlePatchPatient(edit deiz.EditPatient, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var p deiz.Patient
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := validate.StructCtx(ctx, p); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := edit(ctx, &p, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

*/
