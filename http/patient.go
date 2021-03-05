package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	PatientSearcher interface {
		SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error)
	}
	PatientAdder interface {
		AddPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	PatientEditer interface {
		EditPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	PatientAddressAdder interface {
		AddPatientAddress(ctx context.Context, a *deiz.Address, patientID int, clinicianID int) error
	}
	PatientAddressEditer interface {
		EditPatientAddress(ctx context.Context, a *deiz.Address, patientID int, clinicianID int) error
	}
)

func handlePatchPatientAddress(addressEditer PatientAddressEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		patientID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		var a deiz.Address
		if err := c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, errBind)
		}
		if err := addressEditer.EditPatientAddress(ctx, &a, patientID, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostPatientAddress(adder PatientAddressAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		patientID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		var a deiz.Address
		if err := c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, errBind)
		}
		if err := adder.AddPatientAddress(ctx, &a, patientID, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, a)
	}
}

func handlePatchPatient(editer PatientEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var p deiz.Patient
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, errBind)
		}
		if err := editer.EditPatient(ctx, &p, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

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
