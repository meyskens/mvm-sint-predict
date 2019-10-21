package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/meyskens/mvm-sint-predict/pb"
)

type fakeClock struct{}

func (f fakeClock) Since(in time.Time) time.Duration {
	return time.Date(2019, time.October, 19, 0, 0, 0, 0, time.UTC).Sub(in)
}

func Test_childrenCountOptions_readFamilyComposotion(t *testing.T) {
	localClock = fakeClock{} // fake out time
	type fields struct {
		FrequencyFile         string
		FamilyCompositionFile string
		FrequencyThreshold    uint64
	}
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*pb.FamilyComposition
		wantErr bool
	}{
		{
			name: "Run normally",
			args: args{in: strings.NewReader(`id,children
MVM04,2018-03-05 Jongen,2015-11-02 Meisje,2006-08-05 Jongen,2004-06-16 Jongen,2001-04-01 Meisje
MVM02,2011-12-22 Jongen,2007-04-23 Jongen,2005-12-02 Jongen,2001-05-11 Jongen`)},
			wantErr: false,
			want: []*pb.FamilyComposition{
				&pb.FamilyComposition{
					Id: "MVM04",
					Children: []*pb.FamilyComposition_Child{
						&pb.FamilyComposition_Child{
							Age:    1,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    3,
							Gender: "Meisje",
						},
						&pb.FamilyComposition_Child{
							Age:    13,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    15,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    18,
							Gender: "Meisje",
						},
					},
				},
				&pb.FamilyComposition{
					Id: "MVM02",
					Children: []*pb.FamilyComposition_Child{
						&pb.FamilyComposition_Child{
							Age:    7,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    12,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    13,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    18,
							Gender: "Jongen",
						},
					},
				},
			},
		},
		{
			name:    "Run date incorrect",
			args:    args{in: strings.NewReader(`id,children`)},
			wantErr: false,
			want:    []*pb.FamilyComposition{},
		},
		{
			name: "Run normally",
			args: args{in: strings.NewReader(`id,children
MVM04,2018-03-05 Jongen,2015-11-02 Meisje,2006-08-05 Jongen,2004-06-16 Jongen,null-04-01 Meisje`)},
			wantErr: false,
			want: []*pb.FamilyComposition{
				&pb.FamilyComposition{
					Id: "MVM04",
					Children: []*pb.FamilyComposition_Child{
						&pb.FamilyComposition_Child{
							Age:    1,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    3,
							Gender: "Meisje",
						},
						&pb.FamilyComposition_Child{
							Age:    13,
							Gender: "Jongen",
						},
						&pb.FamilyComposition_Child{
							Age:    15,
							Gender: "Jongen",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &childrenCountOptions{
				FrequencyFile:         tt.fields.FrequencyFile,
				FamilyCompositionFile: tt.fields.FamilyCompositionFile,
				FrequencyThreshold:    tt.fields.FrequencyThreshold,
			}
			got, err := c.readFamilyComposotion(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("childrenCountOptions.readFamilyComposotion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("childrenCountOptions.readFamilyComposotion() = %v, want %v", got, tt.want)
			}
		})
	}
}
