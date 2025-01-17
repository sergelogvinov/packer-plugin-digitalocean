package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type Artifact struct {
	// The name of the snapshot
	SnapshotName string

	// The ID of the image
	SnapshotId int

	// The name of the region
	RegionNames []string

	// The client for making API calls
	Client *godo.Client

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}
}

var _ packersdk.Artifact = new(Artifact)

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (*Artifact) Files() []string {
	// No files with DigitalOcean
	return nil
}

func (a *Artifact) Id() string {
	return fmt.Sprintf("%s:%s", strings.Join(a.RegionNames[:], ","), strconv.FormatUint(uint64(a.SnapshotId), 10))
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A snapshot was created: '%v' (ID: %v) in regions '%v'", a.SnapshotName, a.SnapshotId, strings.Join(a.RegionNames[:], ","))
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %d (%s)", a.SnapshotId, a.SnapshotName)
	_, err := a.Client.Images.Delete(context.TODO(), a.SnapshotId)
	return err
}
