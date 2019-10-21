package main

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/meyskens/mvm-sint-predict/pb"
)

func Test_frequencyCmdOptions_readCSV(t *testing.T) {
	type fields struct {
		File string
		Out  string
	}
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*pb.FrequencyRequest_Visit
		wantErr bool
	}{
		{
			name: "Run normally",
			args: args{
				in: strings.NewReader(`date,id
22/08/2019,MVM65
06/08/2019,MVM64`),
			},
			wantErr: false,
			want: []*pb.FrequencyRequest_Visit{
				&pb.FrequencyRequest_Visit{
					Id: "MVM65",
					Date: &pb.Date{
						Day:   22,
						Month: 8,
						Year:  2019,
					},
				},
				&pb.FrequencyRequest_Visit{
					Id: "MVM64",
					Date: &pb.Date{
						Day:   6,
						Month: 8,
						Year:  2019,
					},
				},
			},
		},
		{
			name: "Run empty",
			args: args{
				in: strings.NewReader(`date,id`),
			},
			wantErr: false,
			want:    []*pb.FrequencyRequest_Visit{},
		},
		{
			name: "Run invalid date",
			args: args{
				in: strings.NewReader(`date,id
22-08/2019,MVM65
06/08/2019,MVM64`),
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &frequencyCmdOptions{
				File: tt.fields.File,
				Out:  tt.fields.Out,
			}
			got, err := f.readCSV(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("frequencyCmdOptions.readCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("frequencyCmdOptions.readCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
