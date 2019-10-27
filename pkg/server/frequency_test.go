package server

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/meyskens/mvm-sint-predict/pb"
)

func TestSintReplyServer_GetFrequency(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.FrequencyRequest
	}
	tests := []struct {
		name    string
		s       *SintReplyServer
		args    args
		want    *pb.FrequencyReply
		wantErr bool
	}{
		{
			name: "Test normal simple set",
			args: args{
				ctx: context.Background(),
				in: &pb.FrequencyRequest{
					Visits: []*pb.FrequencyRequest_Visit{
						&pb.FrequencyRequest_Visit{
							Id:   "mvm123",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
						&pb.FrequencyRequest_Visit{
							Id:   "mvm123",
							Date: &pb.Date{Day: 2, Month: 1, Year: 2019},
						},
						&pb.FrequencyRequest_Visit{
							Id:   "mvm456",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
					},
				},
			},
			want: &pb.FrequencyReply{
				Frequencies: []*pb.Frequency{
					&pb.Frequency{
						Id:           "mvm123",
						TimesVisited: 2,
					},
					&pb.Frequency{
						Id:           "mvm456",
						TimesVisited: 1,
					},
				},
			},
		},
		{
			name: "Test normal simple set with duplicates",
			args: args{
				ctx: context.Background(),
				in: &pb.FrequencyRequest{
					Visits: []*pb.FrequencyRequest_Visit{
						&pb.FrequencyRequest_Visit{
							Id:   "mvm123",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
						&pb.FrequencyRequest_Visit{
							Id:   "mvm123",
							Date: &pb.Date{Day: 2, Month: 1, Year: 2019},
						},
						&pb.FrequencyRequest_Visit{
							Id:   "mvm456",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
						&pb.FrequencyRequest_Visit{
							Id:   "mvm456",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
						&pb.FrequencyRequest_Visit{
							Id:   "mvm456",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
					},
				},
			},
			want: &pb.FrequencyReply{
				Frequencies: []*pb.Frequency{
					&pb.Frequency{
						Id:           "mvm123",
						TimesVisited: 2,
					},
					&pb.Frequency{
						Id:           "mvm456",
						TimesVisited: 1,
					},
				},
			},
		},
		{
			name: "Test 1 entry",
			args: args{
				ctx: context.Background(),
				in: &pb.FrequencyRequest{
					Visits: []*pb.FrequencyRequest_Visit{
						&pb.FrequencyRequest_Visit{
							Id:   "mvm123",
							Date: &pb.Date{Day: 1, Month: 1, Year: 2019},
						},
					},
				},
			},
			want: &pb.FrequencyReply{
				Frequencies: []*pb.Frequency{
					&pb.Frequency{
						Id:           "mvm123",
						TimesVisited: 1,
					},
				},
			},
		},
		{
			name: "Test empty",
			args: args{
				ctx: context.Background(),
				in: &pb.FrequencyRequest{
					Visits: []*pb.FrequencyRequest_Visit{},
				},
			},
			want: &pb.FrequencyReply{
				Frequencies: []*pb.Frequency{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SintReplyServer{}
			got, err := s.GetFrequency(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("SintReplyServer.GetFrequency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// sort for DeepEqual
			sort.Slice(got.Frequencies, func(i, j int) bool { return got.Frequencies[i].Id < got.Frequencies[j].Id })
			sort.Slice(tt.want.Frequencies, func(i, j int) bool { return tt.want.Frequencies[i].Id < tt.want.Frequencies[j].Id })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SintReplyServer.GetFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}
