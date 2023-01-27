package handlers

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func Test_logicalExpressionEvaluation(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Test 0",
			args:    args{expression: "(!(1 AND 1) OR !(0 AND 1))"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test 1",
			args:    args{expression: "((1 AND 1) OR 1)"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test 2",
			args:    args{expression: "(!(1 AND 1) OR 0)"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Test 3",
			args:    args{expression: "((0 AND 1) OR 0)"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Test 4",
			args:    args{expression: "(((0 AND 1) OR 0)"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got, err := logicalExpressionEvaluation(tt.args.expression)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("logicalExpressionEvaluation() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if got != tt.want {
			// 	t.Errorf("logicalExpressionEvaluation() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestEvaluateExpression(t *testing.T) {
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
			if err := EvaluateExpression(tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("EvaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
