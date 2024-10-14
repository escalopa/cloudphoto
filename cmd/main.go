package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/escalopa/cloudphoto/internal/s3"
)

var bucket = flag.String("bucket", "cloud-photo-cli-storage", "storage bucket name")

func main() {
	ctx := context.Background()

	s3Client := s3.New()

	f, err := os.Open("Task.pdf")
	if err != nil {
		panic(err)
	}

	err = s3Client.Upload(ctx, *bucket, s3Client.GenKey(), f, "application/pdf")
	if err != nil {
		panic(err)
	}

	folder, err := s3Client.ListObjects(ctx, *bucket, "")
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(folder)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
