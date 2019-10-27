package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"github.com/meyskens/mvm-sint-predict/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const year = time.Hour * 8766

type clockInterface interface {
	Since(in time.Time) time.Duration
}

type clock struct{}

func (c clock) Since(in time.Time) time.Duration {
	return time.Since(in)
}

var localClock clockInterface // needed for running tests properly

func init() {
	rootCmd.AddCommand(NewChildrenCountCmd())
	localClock = clock{}
}

type childrenCountOptions struct {
	FrequencyFile         string
	FamilyCompositionFile string
	FrequencyThreshold    uint64
}

// NewChildrenCountCmd generates the `children-count` command
func NewChildrenCountCmd() *cobra.Command {
	s := childrenCountOptions{}
	c := &cobra.Command{
		Use:   "children-count",
		Short: "Calculates the number of children likely to come",
		Long:  `Calculates the number of children likely to come based on the frequency info and famili composition`,
		RunE:  s.RunE,
	}
	c.Flags().StringVarP(&s.FrequencyFile, "frequency-file", "f", "", "Frequency file to read, output from the frequency command")
	c.Flags().StringVarP(&s.FamilyCompositionFile, "family-composition-file", "c", "", "Family composition CSV file to read")
	c.Flags().Uint64VarP(&s.FrequencyThreshold, "frequency-threshold", "t", 3, "Minimal visits to be considered frequent")
	c.MarkFlagRequired("frequency-file")
	c.MarkFlagRequired("family-composition-file")

	viper.BindPFlags(c.Flags())

	return c
}

func (c *childrenCountOptions) RunE(cmd *cobra.Command, args []string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to connect: %v", err)
	}
	defer conn.Close()

	f, err := os.Open(c.FamilyCompositionFile)
	if err != nil {
		return fmt.Errorf("Could not open file: %v", err)
	}
	compositions, err := c.readFamilyComposotion(f)
	if err != nil {
		return fmt.Errorf("Could parse compositions file: %v", err)
	}

	ff, err := os.Open(c.FrequencyFile)
	if err != nil {
		return fmt.Errorf("Could not open file: %v", err)
	}
	ffData, err := ioutil.ReadAll(ff)
	if err != nil {
		return fmt.Errorf("Could not read frequency file: %v", err)
	}

	frequency := []*pb.Frequency{}
	err = json.Unmarshal(ffData, &frequency)
	if err != nil {
		return fmt.Errorf("Could not unmarshal frequency file: %v", err)
	}

	req := &pb.ChildrenCountRequest{
		FamilyCompositions: compositions,
		Frequency:          frequency,
		FrequencyThreshold: c.FrequencyThreshold,
	}

	client := pb.NewSintClient(conn)
	reply, err := client.GetChildrenCount(context.Background(), req)
	if err != nil {
		return fmt.Errorf("Error from server: %v", err)
	}

	for _, count := range reply.Counts {
		fmt.Printf("%s %d: %d\n", count.GetGender(), count.GetAge(), count.GetCount())
	}
	return nil
}

func (c *childrenCountOptions) readFamilyComposotion(in io.Reader) ([]*pb.FamilyComposition, error) {
	out := []*pb.FamilyComposition{}

	r := csv.NewReader(in)
	r.Read() // skip header
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err, _ := err.(*csv.ParseError); err != nil && err.Err != csv.ErrFieldCount {
			return nil, fmt.Errorf("Error while parsing CSV: %v", err)
		}

		family := &pb.FamilyComposition{
			Id:       record[0],
			Children: []*pb.FamilyComposition_Child{},
		}

		if len(record) > 1 {
			for _, entry := range record[1:] {
				entryParts := strings.Split(entry, " ") // is a hack in the current output of ZOHO, not ideal
				t, err := time.Parse("2006-01-02", entryParts[0])
				if err != nil {
					continue // ignore invalid data for now
				}

				family.Children = append(family.Children, &pb.FamilyComposition_Child{
					Gender: entryParts[1],
					Age:    uint32(math.Floor(float64(localClock.Since(t) / year))),
				})
			}
		}

		out = append(out, family)
	}
	return out, nil
}
