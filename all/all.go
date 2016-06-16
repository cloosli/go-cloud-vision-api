// go run label.go <path-to-image>
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"encoding/base64"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
	"io/ioutil"
)

// run submits a label request on a single image by given file.
func run(file string) error {
	ctx := context.Background()

	// Authenticate to generate a vision service
	client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
	if err != nil {
		return err
	}
	service, err := vision.New(client)
	if err != nil {
		return err
	}

	// Read the image
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Construct a label request, encoding the image in base64
	req := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},
		Features: []*vision.Feature{
			{Type: "TYPE_UNSPECIFIED"},
			{Type: "FACE_DETECTION"},
			{Type: "LANDMARK_DETECTION"},
			{Type: "LOGO_DETECTION"},
			{Type: "LABEL_DETECTION"},
			{Type: "TEXT_DETECTION"},
			{Type: "SAFE_SEARCH_DETECTION"},
			{Type: "IMAGE_PROPERTIES"},
		},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		return err
	}

	p := func(j int, ann *vision.EntityAnnotation) {
		fmt.Printf("\tFound: %d %s \tscore: %v\n", j, ann.Description, ann.Score)
		json,_:=ann.MarshalJSON()
		fmt.Printf("\t\t%s\n", json)
	}
	pf := func(j int, ann *vision.FaceAnnotation) {
		fmt.Printf("\tFound: %d %s \tscore: %v\n", j, ann.HeadwearLikelihood, ann.SurpriseLikelihood)
		json,_:=ann.MarshalJSON()
		fmt.Printf("\t\t%s\n", json)
	}
	pi := func(ann *vision.ImageProperties) {
		fmt.Printf("\tFound: %v\n", ann.DominantColors.Colors[0])
		json,_:=ann.MarshalJSON()
		fmt.Printf("\t\t%s\n", json)
	}

	for _, resp := range res.Responses {
		fmt.Printf("Labels for file: %s\n", file)
		for j, ann := range resp.LabelAnnotations {
			p(j, ann)
		}
		fmt.Printf("Logo for file: %s\n", file)
		for j, ann := range resp.LogoAnnotations {
			p(j, ann)
		}
		fmt.Printf("Text for file: %s\n", file)
		for j, ann := range resp.TextAnnotations {
			p(j, ann)
		}
		fmt.Printf("Face for file: %s\n", file)
		for j, ann := range resp.FaceAnnotations {
			pf(j, ann)
		}
		fmt.Printf("Landmark for file: %s\n", file)
		for j, ann := range resp.LandmarkAnnotations {
			p(j, ann)
		}
		fmt.Printf("Image for file: %s\n", file)
		ann := resp.ImagePropertiesAnnotation
		pi(ann)
	}

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
