package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"

	"github.com/csv-processor/business"
	"github.com/csv-processor/model"
	"github.com/csv-processor/mongo"

	"os"

	"go.uber.org/zap"
)

var (
	processFile string
	parkname    string
	parkslug    string
	parkid      int64
	filetype    string

	log *zap.Logger
)

func init() {
	flag.StringVar(&processFile, "csvFile", "/home/user/file.csv", "path to file with CSV to parse")
	flag.StringVar(&filetype, "filetype", "transactions", "file type to be processed")
	flag.StringVar(&parkname, "parkname", "Monza", "the name of park to get business logic")
	flag.StringVar(&parkslug, "parkslug", "monza", "the slug of park to get business logic")
	flag.Int64Var(&parkid, "parkid", 6, "the id of park to get business logic")

	required := []string{"csvFile", "parkname", "filetype", "parkslug", "parkid"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required [-%s] argument/flag\n", req)
			os.Exit(2)
		}
	}

	log, _ = zap.NewProduction()
}

func main() {
	log.Info("Starting parser")
	defer log.Sync()

	db, err := mongo.NewConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting database: [%s]\n", err.Error())
		os.Exit(2)
	}

	csvIn, err := os.Open(processFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file [%s]: [%s]\n", processFile, err.Error())
		os.Exit(2)
	}
	defer csvIn.Close()
	reader := csv.NewReader(csvIn)

	parking := model.Parking{
		Name: parkname,
		Slug: parkslug,
		ID:   parkid,
	}

	processor := business.NewVP(db, reader, filetype, parking)
	err = processor.Process(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error processing file [%s]: [%s]\n", processFile, err.Error())
		os.Exit(2)
	}

	log.Info("Finishing...")
}
