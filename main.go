package main

import (
	"context"
	"fmt"
	"github.com/xoxoist/cl-tutor/applications/alpha"
	"github.com/xoxoist/cl-tutor/applications/bravo"
	"github.com/xoxoist/cl-tutor/applications/charlie"
	"github.com/xoxoist/cl-tutor/logger"
	"go.uber.org/zap"
)

func ControllerAlphaApplication() (*alpha.Response, error) {
	// create logger
	const appName, logFile, logMode, version = "alpha", "logs/log.svc.a.log", "prd", "v1.0.0"
	lg := logger.NewZapLogger(logFile, logMode, version, appName)

	// prepare contextual logging
	const XRequestIDHeader = "1239893491327"
	requestData := alpha.Request{
		AlphaContent:   "John Doe",
		BravoContent:   "Hello World from Bravo to",
		CharlieContent: "Hello World from Charlie to",
	}
	lg = lg.With("request_id", XRequestIDHeader)
	lg = lg.With("request_data", requestData)

	// prepare context for contextual logging
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Request-Id", XRequestIDHeader)
	ctx = context.WithValue(ctx, zap.SugaredLogger{}, lg)

	// calling business logic
	service := alpha.NewAlphaApplication(
		appName,
		version,
		ControllerBravoApplication,
		ControllerCharlieApplication,
	)
	response, err := service.AlphaLogicFlowEntry(ctx, requestData)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ControllerBravoApplication(requestId string, request bravo.Request) (*bravo.Response, error) {
	// create logger
	const appName, logFile, logMode, version = "bravo", "logs/log.svc.b.log", "prd", "v2.0.0"
	lg := logger.NewZapLogger(logFile, logMode, version, appName)

	// prepare contextual logging
	lg = lg.With("request_id", requestId)
	lg = lg.With("request_data", request)

	// prepare context for contextual logging
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Request-Id", requestId)
	ctx = context.WithValue(ctx, zap.SugaredLogger{}, lg)

	// prepare context for contextual logging
	service := bravo.NewBravoApplication(appName, version)
	return service.BravoLogicFlowEntry(ctx, request)
}

func ControllerCharlieApplication(requestId string, request charlie.Request) (*charlie.Response, error) {
	// create logger
	const appName, logFile, logMode, version = "charlie", "logs/log.svc.c.log", "prd", "v3.0.0"
	lg := logger.NewZapLogger(logFile, logMode, version, appName)

	// prepare contextual logging
	lg = lg.With("request_id", requestId)
	lg = lg.With("request_data", request)

	// prepare context for contextual logging
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Request-Id", requestId)
	ctx = context.WithValue(ctx, zap.SugaredLogger{}, lg)

	// prepare context for contextual logging
	service := charlie.NewCharlieApplication(appName, version)
	return service.CharlieLogicFlowEntry(ctx, request)
}

func main() {
	response, err := ControllerAlphaApplication()
	if err != nil {
		panic(err)
	}
	fmt.Println("BRAVO:", response.BravoManipulatedContent)
	fmt.Println("CHARLIE:", response.CharlieManipulatedContent)
}
