package server

import (
	"context"
	"fmt"

	"github.com/meyskens/mvm-sint-predict/pb"
)

// GetFrequency gives back the frequency a family requests services given raw input data of all visits
func (s *SintReplyServer) GetFrequency(ctx context.Context, in *pb.FrequencyRequest) (*pb.FrequencyReply, error) {
	familyVisists := map[string][]string{}

	for _, visit := range in.Visits {
		if _, exists := familyVisists[visit.GetId()]; !exists {
			familyVisists[visit.GetId()] = []string{}
		}
		familyVisists[visit.GetId()] = append(familyVisists[visit.GetId()], fmt.Sprintf("%d-%d-%d", visit.GetDate().GetYear(), visit.GetDate().GetMonth(), visit.GetDate().GetDay()))
	}

	for key, familyVisist := range familyVisists {
		familyVisists[key] = removeDuplicates(familyVisist)
	}

	reply := &pb.FrequencyReply{
		Frequencies: []*pb.Frequency{},
	}
	for key, familyVisist := range familyVisists {
		reply.Frequencies = append(reply.Frequencies, &pb.Frequency{
			Id:           key,
			TimesVisited: uint64(len(familyVisist)),
		})
	}

	return reply, nil
}

func removeDuplicates(in []string) []string {
	out := []string{}
	hadItem := map[string]bool{}

	for _, item := range in {
		if _, exists := hadItem[item]; !exists {
			hadItem[item] = true
			out = append(out, item)
		}
	}

	return out
}
