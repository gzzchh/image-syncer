package sync

import (
	"fmt"

	"github.com/containers/image/v5/manifest"
	"github.com/opencontainers/go-digest"
)

// ManifestSchemaV2List describes a schema V2 manifest list
type ManifestSchemaV2List struct {
	Manifests []ManifestSchemaV2Info `json:"manifests"`
}

// ManifestSchemaV2Info includes of the imformation needes of a schema V2 manifest file
type ManifestSchemaV2Info struct {
	Digest    string `json:"digest"`
	MediaType string `json:"mediaType"`
	Platform  struct {
		// 补充对于多 Arch 镜像的支持
		Architecture string `json:"architecture"`
		Os           string `json:"os"`
	} `json:"platform"`
	Size int64 `json:"size"`
}

// ManifestHandler expends the ability of handling manifest list in schema2, but it's not finished yet
// return the digest array of manifests in the manifest list if exist.
func ManifestHandler(m []byte, t string, i *ImageSource) (manifest.Manifest, []*digest.Digest, error) {
	if t == manifest.DockerV2Schema2MediaType {
		manifestInfo, err := manifest.Schema2FromManifest(m)
		if err != nil {
			return nil, nil, err
		}
		return manifestInfo, nil, nil
	} else if t == manifest.DockerV2Schema1MediaType {
		manifestInfo, err := manifest.Schema1FromManifest(m)
		if err != nil {
			return nil, nil, err
		}
		return manifestInfo, nil, nil
	} else if t == manifest.DockerV2ListMediaType {
		manifestInfo, err := manifest.Schema2ListFromManifest(m)
		if err != nil {
			return nil, nil, err
		}
		//fmt.Println(manifestInfo.Manifests)
		//manifestDigests := []*digest.Digest{}
		for _, item := range manifestInfo.Manifests {
			// 暂时至支持X86_64架构的镜像同步
			if item.Platform.Architecture == "amd64" {
				// For Linux Only
				if item.Platform.OS == "linux" {
					//fmt.Println("newCtx:", i.ctx)
					//fmt.Println("newDigest:", item.Digest)
					// 用指定Arch的Digest再拉一次针对该Digest的多层Manifest
					manifestByte, manifestType, err := i.source.GetManifest(i.ctx, &item.Digest)
					if err != nil {
						return nil, nil, err
					}
					platformSpecManifest, _, err := ManifestHandler(manifestByte, manifestType, i)
					//fmt.Println(platformSpecManifest)
					return platformSpecManifest, nil, nil
				}

			}
		}
		//fmt.Println(manifestInfo.Manifests)
		//return nil, manifestDigests, nil
		return nil, nil, nil
	}
	return nil, nil, fmt.Errorf("unsupported manifest type: %v", t)
}
