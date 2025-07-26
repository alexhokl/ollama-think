package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func run(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: %s <model> <prompt>", args[0])
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	shouldThink := true
	req := &api.GenerateRequest{
		Model:  args[1],
		Prompt: args[2],
		Think:  &shouldThink,
	}

	// create to memorry streams to capture the response and thinking output
	normalOutput := strings.Builder{}
	thinkingOutput := strings.Builder{}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		normalOutput.WriteString(resp.Response)
		thinkingOutput.WriteString(resp.Thinking)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("Response:", normalOutput.String())

	if thinkingOutput.Len() > 0 {
		fmt.Println("Thinking:", thinkingOutput.String())
	}

	return nil
}
