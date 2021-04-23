package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handlePatchPatientAddress(addressEditer usecase.PatientAddressEditer) echo.HandlerFunc {
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

func handlePostPatientAddress(adder usecase.PatientAddressAdder) echo.HandlerFunc {
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

func handlePatchPatient(editer usecase.PatientEditer) echo.HandlerFunc {
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

func handleGetPatients(searcher usecase.PatientSearcher) echo.HandlerFunc {
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

func handlePostPatient(adder usecase.PatientAdder) echo.HandlerFunc {
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
