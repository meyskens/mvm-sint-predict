package server

import (
	"context"

	"github.com/meyskens/mvm-sint-predict/pb"
)

type childType struct {
	Age    uint32
	Gender string
}

// GetChildrenCount tries to predict how many children will visit based on frequency data and family composition
func (s *SintReplyServer) GetChildrenCount(ctx context.Context, in *pb.ChildrenCountRequest) (*pb.ChildrenCountReply, error) {
	frequentFamilies := map[string]bool{}
	childCount := map[childType]*pb.ChildrenCountReply_Count{}

	for _, frequency := range in.GetFrequency() {
		if frequency.GetTimesVisited() >= in.GetFrequencyThreshold() {
			frequentFamilies[frequency.GetId()] = true
		}
	}

	for _, familyComposition := range in.GetFamilyCompositions() {
		if _, isFrequent := frequentFamilies[familyComposition.GetId()]; isFrequent {
			for _, child := range familyComposition.GetChildren() {
				child.GetAge()
				child.GetGender()
				t := childType{
					Age:    child.GetAge(),
					Gender: child.GetGender(),
				}
				if _, exists := childCount[t]; !exists {
					childCount[t] = &pb.ChildrenCountReply_Count{
						Age:    child.GetAge(),
						Gender: child.GetGender(),
						Count:  0,
					}
				}

				childCount[t].Count++
			}
		}
	}

	resp := &pb.ChildrenCountReply{
		Counts: []*pb.ChildrenCountReply_Count{},
	}

	for _, childReplyCount := range childCount {
		resp.Counts = append(resp.Counts, childReplyCount)
	}

	return resp, nil
}
