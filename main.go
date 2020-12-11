package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/soap-parser/business"
	"github.com/soap-parser/model"
	"github.com/soap-parser/mongo"
	"io/ioutil"

	"go.uber.org/zap"
	"os"
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
	flag.StringVar(&processFile, "xmlFilePath", "/home/user/request.xml", "path to file with XML to parse")
	flag.StringVar(&filetype, "filetype", "pagamentos", "file type to be processed")
	flag.StringVar(&parkname, "parkname", "Monza", "the name of park to get business logic")
	flag.StringVar(&parkslug, "parkslug", "monza", "the slug of park to get business logic")
	flag.Int64Var(&parkid, "parkid", 6, "the id of park to get business logic")

	required := []string{"xmlFilePath", "parkname", "filetype", "parkslug", "parkid"}
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
	log.Info("Starting soap parser")
	defer log.Sync()

	db, err := mongo.NewConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting database: [%s]\n", err.Error())
		os.Exit(2)
	}

	file, err := ioutil.ReadFile(processFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file [%s]: [%s]\n", processFile, err.Error())
		os.Exit(2)
	}

	parking := model.Parking{
		Name: parkname,
		Slug: parkslug,
		ID:   parkid,
	}

	processor := business.NewAuconMonza(db, file, filetype, parking)
	err = processor.Process(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error processing file [%s]: [%s]\n", processFile, err.Error())
		os.Exit(2)
	}

	log.Info("Finishing...")
}
