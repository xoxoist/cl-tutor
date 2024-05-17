package alpha

import (
	"context"
	"errors"
	"fmt"
	"github.com/xoxoist/cl-tutor/applications/bravo"
	"github.com/xoxoist/cl-tutor/applications/charlie"
	"go.uber.org/zap"
)

type (
	Request struct {
		AlphaContent   string
		BravoContent   string
		CharlieContent string
	}

	Response struct {
		BravoManipulatedContent   string
		CharlieManipulatedContent string
	}

	Alpha interface {
		AlphaLogicFlowEntry(ctx context.Context, req Request) (*Response, error)
		businessLogicFlow(ctx context.Context, req Request) (*Response, error)
		requestToBravoApp(ctx context.Context, req Request) (*bravo.Response, error)
		requestToCharlieApp(ctx context.Context, req Request) (*charlie.Response, error)
	}

	alpha struct {
		Name          string
		Version       string
		Package       string
		bravoClient   func(requestId string, request bravo.Request) (*bravo.Response, error)
		charlieClient func(requestId string, request charlie.Request) (*charlie.Response, error)
	}
)

func NewAlphaApplication(
	applicationName string,
	applicationVersion string,
	bravoClient func(requestId string, request bravo.Request) (*bravo.Response, error),
	charlieClient func(requestId string, request charlie.Request) (*charlie.Response, error),
) Alpha {
	return &alpha{
		Name:          applicationName,
		Version:       applicationVersion,
		Package:       "applications.alpha",
		bravoClient:   bravoClient,
		charlieClient: charlieClient,
	}
}

func (a *alpha) AlphaLogicFlowEntry(ctx context.Context, req Request) (*Response, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", a.Package, "AlphaLogicFlowEntry")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return nil, errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// executing business logic flow
	response, err := a.businessLogicFlow(ctx, req)
	if err != nil {
		cl.Errorw(fName, "reason", err.Error())
		return nil, errors.New("error while contacting charlie service")
	}

	// returning
	return response, nil
}

func (a *alpha) businessLogicFlow(ctx context.Context, req Request) (*Response, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", a.Package, "businessLogicFlow")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return nil, errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// var bravo and charlie manipulated string
	var bravoStr, charlieStr string

	// when string length inside req.BravoContent is 0
	// then calling bravo client
	if len(req.BravoContent) > 0 {
		br, err := a.requestToBravoApp(ctx, req)
		if err != nil {
			cl.Errorw(fName, "reason", err.Error())
			return nil, errors.New("error while contacting bravo service")
		}
		bravoStr = br.ManipulatedContent
	}

	// when string length inside req.CharlieContent is 0
	// then calling charlie client
	if len(req.CharlieContent) > 0 {
		cr, err := a.requestToCharlieApp(ctx, req)
		if err != nil {
			cl.Errorw(fName, "reason", err.Error())
			return nil, errors.New("error while contacting charlie service")
		}
		charlieStr = cr.ManipulatedContent
	}

	// returning
	return &Response{
		BravoManipulatedContent:   fmt.Sprintf("%s %s", req.AlphaContent, bravoStr),
		CharlieManipulatedContent: fmt.Sprintf("%s %s", req.AlphaContent, charlieStr),
	}, nil
}

func (a *alpha) requestToBravoApp(ctx context.Context, req Request) (*bravo.Response, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", a.Package, "requestToBravoApp")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return nil, errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// extract request id to be passed to bravo application
	requestId, ok := ctx.Value("Request-Id").(string)
	if !ok {
		return nil, errors.New("request id not pass through context")
	}

	// sending request to bravo application
	response, err := a.bravoClient(requestId, bravo.Request{
		Content: req.BravoContent,
	})
	if err != nil {
		cl.Errorw(fName, "reason", err.Error())
		return nil, errors.New("error while contacting bravo service")
	}

	// returning
	return response, nil
}

func (a *alpha) requestToCharlieApp(ctx context.Context, req Request) (*charlie.Response, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", a.Package, "requestToCharlieApp")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return nil, errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// extract request id to be passed to charlie application
	requestId, ok := ctx.Value("Request-Id").(string)
	if !ok {
		return nil, errors.New("request id not pass through context")
	}

	// sending request to charlie application
	response, err := a.charlieClient(requestId, charlie.Request{
		Content: req.CharlieContent,
	})
	if err != nil {
		cl.Errorw(fName, "reason", err.Error())
		return nil, errors.New("error while contacting charlie service")
	}

	// returning
	return response, nil
}
