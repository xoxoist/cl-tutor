package charlie

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

	Charlie interface {
		CharlieLogicFlowEntry(ctx context.Context, req Request) (*Response, error)
		replaceVowelsWithAsterisk(ctx context.Context, content string) (string, error)
	}

	charlie struct {
		Name    string
		Version string
		Package string
	}
)

func NewCharlieApplication(
	applicationName string,
	applicationVersion string,
) Charlie {
	return &charlie{
		Name:    applicationName,
		Version: applicationVersion,
		Package: "applications.charlie",
	}
}

func (c *charlie) CharlieLogicFlowEntry(ctx context.Context, req Request) (*Response, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", c.Package, "BravoLogicFlowEntry")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return nil, errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// reverse the string
	manipulated, err := c.replaceVowelsWithAsterisk(ctx, req.Content)
	if err != nil {
		cl.Errorw(fName, "reason", err.Error())
		return nil, errors.New("error while manipulate string")
	}

	// returning
	return &Response{
		ManipulatedContent: manipulated,
	}, nil
}

func (c *charlie) replaceVowelsWithAsterisk(ctx context.Context, content string) (string, error) {
	// extract zap logger inside context that passed from controller where
	// request received and logger from controller its already has unique
	// request id everytime request received.
	fName := fmt.Sprintf("%s.%s", c.Package, "replaceVowelsWithAsterisk")
	cl, ok := ctx.Value(zap.SugaredLogger{}).(*zap.SugaredLogger)
	if !ok {
		return "", errors.New("contextual logging not pass through context")
	}
	cl.Infow(fName, "reason", "execution started")
	defer cl.Infow(fName, "reason", "execution ended")

	// manipulator
	vowels := "AEIOUaeiou"
	var result strings.Builder
	for _, char := range content {
		if strings.ContainsRune(vowels, char) {
			result.WriteRune('*')
		} else {
			result.WriteRune(char)
		}
	}
	manipulated := result.String()
	cl.Infow(fName, "manipulated_content", manipulated)

	// returning
	return manipulated, nil
}
