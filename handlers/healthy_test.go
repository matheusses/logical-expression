package handlers

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthy(t *testing.T) {
	type args struct {
		context *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Healthy(tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("Healthy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
