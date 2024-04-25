package handlers

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/amitrei/oas-compare-server/database"
	"github.com/amitrei/oas-compare-server/logger"
	"github.com/labstack/echo"
	"github.com/pb33f/openapi-changes/cmd"
	"github.com/pb33f/openapi-changes/model"
	"github.com/twinj/uuid"
	"gopkg.in/validator.v2"
)

//go:embed  resources/notfound-page.html
var notFoundPage string

type HttpHandler struct {
	HandlerFunc echo.HandlerFunc
	Method      string
	Path        string
}

type createReportRequest struct {
	Original string `json:"original" validate:"nonzero,nonnil"`
	Modified string `json:"modified" validate:"nonzero,nonnil"`
}

type createReportResponse struct {
	ReportId string `json:"reportId"`
	Error    string `json:"error"`
}

func GlobalErrorHandler(error error, ctx echo.Context) {
	if error != nil {
		logger := logger.GetContextLogger(ctx)
		errMessage := fmt.Sprintf("An error was thrown: %s", error.Error())
		logger.Error(errMessage)
		ctx.JSON(400, map[string]string{"error": error.Error()})
	}
}

func GetHandlers() []HttpHandler {
	return []HttpHandler{CreateReportHandler(), GetReportHandler()}
}

func GetReportHandler() HttpHandler {
	handlerFunc := func(ctx echo.Context) error {
		id := ctx.Param("id")
		report, err := database.GetDatabaseClient().Get(id)
		if err != nil {
			return ctx.HTML(200, notFoundPage)
		}
		return ctx.HTML(200, report.(string))
	}

	return HttpHandler{HandlerFunc: handlerFunc, Method: http.MethodGet, Path: "/report/:id"}
}

func CreateReportHandler() HttpHandler {
	handlerFunc := func(ctx echo.Context) error {
		reportId := uuid.NewV4().String()

		requestBody := createReportRequest{}
		_ = ctx.Bind(&requestBody)

		report, err := generateReport(&requestBody)
		if err != nil {
			return err
		}

		err = database.GetDatabaseClient().Set(reportId, report)

		if err != nil {
			return err
		}

		return ctx.JSON(200, createReportResponse{ReportId: reportId})
	}

	return HttpHandler{HandlerFunc: handlerFunc, Method: http.MethodPost, Path: "/compare"}
}

func generateReport(request *createReportRequest) ([]byte, error) {
	err := validator.Validate(request)
	if err != nil {
		return nil, err
	}

	// Wrapping compare function in a go routine due to an internal Fatal throwing that cannot be handled In cases of schema references not found.
	doneChan := make(chan []byte, 1)
	errChan := make(chan model.ProgressError, 100)
	progressChan := make(chan *model.ProgressUpdate, 100)

	go func() {
		report, _ := cmd.RunLeftRightHTMLReportViaString(string(request.Original),
			string(request.Modified), false, false,
			progressChan, errChan, "", false)
		doneChan <- report
	}()

	for {
		select {
		case err := <-errChan:
			close(doneChan)
			return nil, err

		case report := <-doneChan:
			close(doneChan)
			return report, nil
		}
	}

}
