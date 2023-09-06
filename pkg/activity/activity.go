package activity

import (
	"context"
	"fmt"
)

type Activities struct{}

func (a *Activities) SampleActivity1(ctx context.Context, input []string) (string, error) {
	fmt.Printf("Run with input %v \n", input)
	return "Result_", nil
}

func (a *Activities) SampleActivity2(ctx context.Context, input []string) (string, error) {
	fmt.Printf("Run with input %v \n", input)
	return "Result_", nil
}

func (a *Activities) SampleActivity3(ctx context.Context, input []string) (string, error) {
	fmt.Printf("Run with input %v \n", input)
	return "Result_", nil
}

func (a *Activities) SampleActivity4(ctx context.Context, input []string) (string, error) {
	fmt.Printf("Run with input %v \n", input)
	return "Result_", nil
}
