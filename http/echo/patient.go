package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

func handlePatchPatientAddress(addressEditer usecase.PatientAddressEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		patientID, err := getURLIntegerParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		a, err := getAddressFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
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
		clinicianID := getCredFromEchoCtx(c).UserID
		patientID, err := getURLIntegerParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		a, err := getAddressFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
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
		clinicianID := getCredFromEchoCtx(c).UserID
		p, err := getPatientFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := editer.EditPatient(ctx, &p, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func getPatientFromRequest(c echo.Context) (deiz.Patient, error) {
	var p deiz.Patient
	return p, c.Bind(&p)
}

func handleGetPatients(searcher usecase.PatientSearcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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
		clinicianID := getCredFromEchoCtx(c).UserID
		p, err := getPatientFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errBind)
		}
		if err := adder.AddPatient(ctx, &p, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, p)
	}
}
