package bravo

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type (
	Request struct {
		Content string
	}

	Response struct {
		ManipulatedContent string
	}

	Bravo interface {
		BravoLogicFlowEntry(ctx context.Context, req Request) (*Response, error)
		reverseString(ctx context.Context, content string) (string, error)
	}

	bravo struct {
		Name    string
		Version string
		Package string
	}
)

func NewBravoApplication(
	applicationName string,
	applicationVersion string,
) Bravo {
	return &bravo{
		Name:    applicationName,
		Version: applicationVersion,
		Package: "applications.bravo",
	}
}

func (b *bravo) BravoLogicFlowEntry(ctx context.Context, req Request) (*Response, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", b.Package, "BravoLogicFlowEntry")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return nil, errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// reverse the string
	manipulated, err := b.reverseString(ctx, req.Content)
	if err != nil {
		cl.Errorw(fName, "reason", err.Error())
		return nil, errors.New("error while manipulate string")
	}

	// returning
	return &Response{
		ManipulatedContent: manipulated,
	}, nil
}

func (b *bravo) reverseString(ctx context.Context, content string) (string, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", b.Package, "reverseString")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return "", errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// manipulator
	runes := []rune(content)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	manipulated := strings.ToUpper(string(runes))
	cl.Infow(fName, "manipulated_content", manipulated)

	// returning
	return manipulated, nil
}
