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
			name: "(x OR y) AND z => x=1&y=1&z=1: true",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "(x OR y) AND z",
			},
			args: args{
				queryString: "x=1&y=1&z=1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "(x OR y) AND z => x=1&y=0&z=1: true",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "(x OR y) AND z",
			},
			args: args{
				queryString: "x=1&y=0&z=1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "(x OR y) AND z => x=1&y=1&z=0: false",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "(x OR y) AND z",
			},
			args: args{
				queryString: "x=1&y=1&z=0",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "((x OR y) AND (z OR k) OR j) => x=1&y=1&z=1&k=1&j=1 :true",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "((x OR y) AND (z OR k) OR j)",
			},
			args: args{
				queryString: "x=1&y=1&z=1&k=1&j=1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "((x OR y) AND (z OR k) OR j) => x=0&y=0&z=1&k=1&j=1 :false",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "((x OR y) AND (z OR k) OR j)",
			},
			args: args{
				queryString: "x=0&y=0&z=1&k=1&j=1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "((x OR y) AND (z and k) AND j) => x=1&y=1&z=1&k=1&j=1 :true",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "((x OR y) AND (z OR k) OR j)",
			},
			args: args{
				queryString: "x=1&y=1&z=1&k=1&j=1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "((x OR y) AND (z and k) AND j) => x=1&y=1&z=1&k=1&j=0 :false",
			fields: fields{
				Model:      gorm.Model{},
				Expression: "((x OR y) AND (z OR k) OR j)",
			},
			args: args{
				queryString: "x=1&y=1&z=1&k=1&j=0",
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
