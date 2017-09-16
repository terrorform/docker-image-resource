package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/transport"
)

type inRequest struct {
	Source struct {
		Repo string `json:"repository"`
		Tag  string `json:"tag"`
	} `json:"source"`
	Params struct {
		RootFS       bool `json:"rootfs"`
		SkipDownload bool `json:"skip_download"`
		Save         bool `json:"save"`
	} `json:"params"`
	Version struct {
		Digest string `json:"digest"`
	} `json:"version"`
}

type inResponse struct {
	Version struct {
		Version string `json:"version"`
	} `json:"version"`
	Metadata []metadataField `json:"metadata"`
}

type metadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const officialRegistry = "registry-1.docker.io"

func main() {
	var in inRequest

	err := json.NewDecoder(os.Stdin).Decode(&in)
	if err != nil {
		panic(err)
	}

	assetsDir := os.Args[1]

	trans := transport.NewTransport(http.DefaultTransport, nil)

	name, err := reference.WithName(in.Source.Repo)
	if err != nil {
		panic(err)
	}

	d := digest.ParseDigest(in.Version.Digest)
	if err != nil {
		panic(err)
	}

	ref, err := reference.WithDigest(name, d)
	if err != nil {
		panic(err)
	}

	repo := client.NewRepository(context.TODO(), ref, officialRegistry, trans)

	tags, err := repo.Tags(context.TODO())
	if err != nil {
		panic(err)
	}
}
