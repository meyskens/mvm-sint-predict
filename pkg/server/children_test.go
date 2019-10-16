package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/meyskens/mvm-sint-predict/pb"
)

func TestSintReplyServer_GetChildrenCount(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.ChildrenCountRequest
	}
	tests := []struct {
		name    string
		s       *SintReplyServer
		args    args
		want    *pb.ChildrenCountReply
		wantErr bool
	}{
		{
			name: "Test normal 2 families, 1 frequent",
			args: args{
				ctx: context.Background(),
				in: &pb.ChildrenCountRequest{
					FrequencyThreshold: 10,
					Frequency: []*pb.Frequency{
						&pb.Frequency{
							Id:           "mvm123",
							TimesVisited: 20,
						},
						&pb.Frequency{
							Id:           "mvm456",
							TimesVisited: 2,
						},
					},
					FamilyCompositions: []*pb.FamilyComposition{
						&pb.FamilyComposition{
							Id: "mvm123",
							Children: []*pb.FamilyComposition_Child{
								&pb.FamilyComposition_Child{
									Age:    5,
									Gender: "boy",
								},
							},
						},
						&pb.FamilyComposition{
							Id: "mvm456",
							Children: []*pb.FamilyComposition_Child{
								&pb.FamilyComposition_Child{
									Age:    10,
									Gender: "girl",
								},
							},
						},
					},
				},
			},
			want: &pb.ChildrenCountReply{
				Counts: []*pb.ChildrenCountReply_Count{
					&pb.ChildrenCountReply_Count{
						Age:    5,
						Gender: "boy",
						Count:  1,
					},
				},
			},
		},
		{
			name: "Test normal 2 families, 0 frequent",
			args: args{
				ctx: context.Background(),
				in: &pb.ChildrenCountRequest{
					FrequencyThreshold: 10,
					Frequency: []*pb.Frequency{
						&pb.Frequency{
							Id:           "mvm123",
							TimesVisited: 2,
						},
						&pb.Frequency{
							Id:           "mvm456",
							TimesVisited: 2,
						},
					},
					FamilyCompositions: []*pb.FamilyComposition{
						&pb.FamilyComposition{
							Id: "mvm123",
							Children: []*pb.FamilyComposition_Child{
								&pb.FamilyComposition_Child{
									Age:    5,
									Gender: "boy",
								},
							},
						},
						&pb.FamilyComposition{
							Id: "mvm456",
							Children: []*pb.FamilyComposition_Child{
								&pb.FamilyComposition_Child{
									Age:    10,
									Gender: "girl",
								},
							},
						},
					},
				},
			},
			want: &pb.ChildrenCountReply{
				Counts: []*pb.ChildrenCountReply_Count{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SintReplyServer{}
			got, err := s.GetChildrenCount(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("SintReplyServer.GetChildrenCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SintReplyServer.GetChildrenCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
