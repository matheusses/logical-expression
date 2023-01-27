package models

import (
	"testing"

	"gorm.io/gorm"
)

func TestLogicalExpression_EvaluatePerQueryString(t *testing.T) {
	type fields struct {
		Model      gorm.Model
		Expression string
	}
	type args struct {
		queryString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "simple test",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "((x AND y) OR z)",
			},
			args: args{
				queryString: "x=1&y=1&z=1",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &LogicalExpression{
				Model:      tt.fields.Model,
				Expression: tt.fields.Expression,
			}
			got, err := self.EvaluatePerQueryString(tt.args.queryString)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicalExpression.EvaluatePerQueryString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LogicalExpression.EvaluatePerQueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}
