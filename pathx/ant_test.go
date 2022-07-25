package pathx

import "testing"

func TestAntMatch(t *testing.T) {
	type args struct {
		path     string
		patterns []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "should match one pattern",
			args: args{
				path:     "testdata",
				patterns: []string{"testdata/**"},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AntMatch(tt.args.path, tt.args.patterns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("AntMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AntMatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
