package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/meyskens/mvm-sint-predict/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(NewFrequencyCmd())
}

type frequencyCmdOptions struct {
	File string
	Out  string
}

// NewFrequencyCmd generates the `frequency` command
func NewFrequencyCmd() *cobra.Command {
	f := frequencyCmdOptions{}
	c := &cobra.Command{
		Use:   "frequency",
		Short: "Calculates the frequency info",
		Long:  `Calculates the frequency info for a given CSV and prints in in JSON`,
		RunE:  f.RunE,
	}
	c.Flags().StringVarP(&f.File, "file", "f", "", "File to read")
	c.Flags().StringVarP(&f.Out, "out", "o", "", "File to write output to")
	c.MarkFlagRequired("file")

	viper.BindPFlags(c.Flags())

	return c
}

func (f *frequencyCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to connect: %v", err)
	}
	defer conn.Close()

	file, err := os.Open(f.File)
	if err != nil {
		return fmt.Errorf("Could not open file: %v", err)
	}
	visits, err := f.readCSV(file)
	if err != nil {
		return fmt.Errorf("Could not parse file: %v", err)
	}

	req := &pb.FrequencyRequest{
		Visits: visits,
	}

	c := pb.NewSintClient(conn)
	reply, err := c.GetFrequency(context.Background(), req)
	if err != nil {
		return fmt.Errorf("Error from server: %v", err)
	}

	out, err := json.Marshal(reply.GetFrequencies())
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %v", err)
	}

	if f.Out == "" {
		fmt.Println(string(out))
		return nil
	}

	err = ioutil.WriteFile(f.Out, out, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON: %v", err)
	}

	return nil
}

func (f *frequencyCmdOptions) readCSV(in io.Reader) ([]*pb.FrequencyRequest_Visit, error) {
	out := []*pb.FrequencyRequest_Visit{}

	r := csv.NewReader(in)
	r.Read() // skip header
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error while parsing CSV: %v", err)
		}

		t, err := time.Parse("02/01/2006", record[0])
		if err != nil {
			return nil, fmt.Errorf("Error while parsing date %s: %v", record[0], err)
		}

		out = append(out, &pb.FrequencyRequest_Visit{
			Date: &pb.Date{
				Day:   int32(t.Day()),
				Month: int32(t.Month()),
				Year:  int32(t.Year()),
			},
			Id: record[1],
		})
	}

	return out, nil
}
