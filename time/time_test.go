package time

import "testing"

func TestLocalTimeFormatFromUnix(t *testing.T) {
	type args struct {
		unixTime int
		format   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Без формата",
			args: args{
				unixTime: 1640445330,
				format:   []string{},
			},
			want: "25.12.2021 18:15:30",
		},
		{
			name: "С форматом RFC3339",
			args: args{
				unixTime: 1640445330,
				format:   []string{"2006-01-02T15:04:05Z07:00"},
			},
			want: "2021-12-25T18:15:30+03:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LocalTimeFormatFromUnix(tt.args.unixTime, tt.args.format...); got != tt.want {
				t.Errorf("LocalTimeFormatFromUnix() = %v, want %v", got, tt.want)
			}
		})
	}
}
