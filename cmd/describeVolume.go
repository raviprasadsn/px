// Copyright © 2019 Portworx
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"text/tabwriter"

	"github.com/cheynewallace/tabby"
	humanize "github.com/dustin/go-humanize"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	prototime "github.com/portworx/px/pkg/openstorage/proto/time"
	"github.com/portworx/px/pkg/portworx"
	"github.com/portworx/px/pkg/util"
	"github.com/spf13/cobra"
)

const (
	timeLayout = "Jan 2 15:04:05 UTC 2006"
)

// describeVolumeCmd represents the describeVolume command
var describeVolumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"volumes"},
	Short:   "Describe a Portworx volume",
	Long:    "Show detailed information of Portworx volumes",
	Example: `$ px describe volume
  This describes all volumes
$ px describe volume abc
  This describes volume abc
$ px describe volume abc xyz
  This describes volumes abc and xyz`,
	RunE: describeVolumesExec,
}

func init() {
	describeCmd.AddCommand(describeVolumeCmd)
	describeVolumeCmd.Flags().String("owner", "", "Owner of volume")
	describeVolumeCmd.Flags().String("volumegroup", "", "Volume group id")
	describeVolumeCmd.Flags().Bool("deep", false, "Collect more information, this may delay the request")
	describeVolumeCmd.Flags().Bool("show-k8s-info", false, "Show kubernetes information")
}

func describeVolumesExec(cmd *cobra.Command, args []string) error {
	vf, err := newVolumeFormatter(cmd, args)
	if err != nil {
		return err
	}
	defer vf.close()

	vcf := volumeInspectFormatter{
		volumeFormatter: *vf,
	}
	vcf.Print()
	return nil
}

type volumeInspectFormatter struct {
	volumeFormatter
}

// String returns the formatted output of the object as per the format set.
func (p *volumeInspectFormatter) String() string {
	return util.GetFormattedOutput(p)
}

// Print writes the object to stdout
func (p *volumeInspectFormatter) Print() {
	util.Printf("%v\n", p)
}

// YamlFormat returns the default representation as there is no yaml format support for describe
func (p *volumeInspectFormatter) YamlFormat() string {
	return p.DefaultFormat()
}

// JsonFormat returns the default representation as there is no json format support for describe
func (p *volumeInspectFormatter) JsonFormat() string {
	return p.DefaultFormat()
}

// WideFormat returns the default representation as there is no wide format support for describe
func (p *volumeInspectFormatter) WideFormat() string {
	return p.DefaultFormat()
}

// DefaultFormat returns the default string representation of the object
func (p *volumeInspectFormatter) DefaultFormat() string {
	return p.toTabbed()
}

func (p *volumeInspectFormatter) toTabbed() string {
	var b bytes.Buffer
	writer := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	t := tabby.NewCustom(writer)

	vols, err := p.pxVolumeOps.GetVolumes()
	if err != nil {
		util.Eprintf("%v\n", err)
		return ""
	}

	for i, n := range vols {
		p.addVolumeDetails(n, t)
		// Put two empty lines between volumes
		if len(vols) > 1 && i != len(vols)-1 {
			t.AddLine("")
			t.AddLine("")
		}
	}
	t.Print()

	return b.String()
}

func (p *volumeInspectFormatter) addVolumeDetails(
	resp *api.SdkVolumeInspectResponse,
	t *tabby.Tabby,
) {

	v := resp.GetVolume()
	p.addVolumeBasicInfo(v, t)
	p.addVolumeStatsInfo(v, t)
	p.addVolumeReplicationInfo(v, t)
	p.addVolumeK8sInfo(v, t)

}

func (p *volumeInspectFormatter) addVolumeBasicInfo(
	v *api.Volume,
	t *tabby.Tabby,
) {
	spec := v.GetSpec()

	// Determine the state of the volume
	state := p.pxVolumeOps.GetAttachedState(v)

	// Print basic info
	t.AddLine("Volume:", v.GetId())
	t.AddLine("Name:", v.GetLocator().GetName())
	if v.GetGroup() != nil && len(v.GetGroup().GetId()) != 0 {
		t.AddLine("Group:", v.GetGroup().GetId())
	}
	if v.GetFormat() == api.FSType_FS_TYPE_FUSE {
		t.AddLine("Type:", "Namespace Volume Group")
		return
	}
	t.AddLine("Size:", humanize.BigIBytes(big.NewInt(int64(spec.GetSize()))))
	t.AddLine("Format:",
		strings.TrimPrefix(v.GetFormat().String(), "FS_TYPE_"))
	t.AddLine("HA:", spec.GetHaLevel())
	t.AddLine("IO Priority:", spec.GetCos())
	t.AddLine("Creation Time:",
		prototime.TimestampToTime(v.GetCtime()).Format(timeLayout))
	if v.GetSource() != nil && len(v.GetSource().GetParent()) != 0 {
		t.AddLine("Parent:", v.GetSource().GetParent())
	}
	snapSched := portworx.SchedSummary(v)
	if len(snapSched) != 0 {
		util.AddArray(t, "Snapshot Schedule:", snapSched)
	}
	if spec.GetStoragePolicy() != "" {
		t.AddLine("StoragePolicy:", spec.GetStoragePolicy())
	}
	t.AddLine("Shared:", portworx.SharedString(v))
	t.AddLine("Status:", portworx.PrettyStatus(v))
	t.AddLine("State:", state)
	attrs := portworx.BooleanAttributes(v)
	if len(attrs) != 0 {
		util.AddArray(t, "Attributes:", attrs)
	}
	if spec.GetScale() > 1 {
		t.AddLine("Scale:", v.Spec.Scale)
	}
	if v.GetAttachedOn() != "" && v.GetAttachedState() != api.AttachState_ATTACH_STATE_INTERNAL {
		t.AddLine("Device Path:", v.GetDevicePath())
	}
	if len(v.GetLocator().GetVolumeLabels()) != 0 {
		util.AddMap(t, "Labels:", v.GetLocator().GetVolumeLabels())
	}
}

func (p *volumeInspectFormatter) addVolumeStatsInfo(
	v *api.Volume,
	t *tabby.Tabby,
) {
	stats := p.pxVolumeOps.GetStats(v)
	t.AddLine("Stats:")
	t.AddLine("  Reads:", stats.GetReads())
	t.AddLine("  Reads MS:", stats.GetReadMs())
	t.AddLine("  Bytes Read:", stats.GetReadBytes())
	t.AddLine("  Writes:", stats.GetWrites())
	t.AddLine("  Writes MS:", stats.GetWriteMs())
	t.AddLine("  Bytes Written:", stats.GetWriteBytes())
	t.AddLine("  IOs in progress:", stats.GetIoProgress())
	t.AddLine("  Bytes used:", humanize.BigIBytes(big.NewInt(int64(stats.BytesUsed))))
}

func (p *volumeInspectFormatter) addVolumeReplicationInfo(
	v *api.Volume,
	t *tabby.Tabby,
) {
	replInfo := p.pxVolumeOps.GetReplicationInfo(v)
	t.AddLine("Replication Status:", replInfo.Status)
	if len(replInfo.Rsi) > 0 {
		t.AddLine("Replica sets on nodes:")
	}
	for _, rsi := range replInfo.Rsi {
		t.AddLine("  Set:", rsi.Id)
		util.AddArray(t, "    Node:", rsi.NodeInfo)
		if len(rsi.HaIncrease) > 0 {
			t.AddLine("    HA-Increase on:", rsi.HaIncrease)
		}
		if len(rsi.ReAddOn) > 0 {
			util.AddArray(t, "    Re-add on:", rsi.ReAddOn)
		}
	}
}

func (p *volumeInspectFormatter) addVolumeK8sInfo(
	v *api.Volume,
	t *tabby.Tabby,
) {
	usedPods := p.pxVolumeOps.PodsUsingVolume(v)
	if len(usedPods) > 0 {
		t.AddLine("Pods:")
		for _, consumer := range usedPods {
			t.AddLine("  - Name:", fmt.Sprintf("%s (%s)",
				consumer.GetName(), consumer.GetUID()))
			t.AddLine("    Namespace:", consumer.GetNamespace())
			t.AddLine("    Running on:", consumer.Spec.NodeName)
			o := make([]string, 0)
			for _, owner := range consumer.OwnerReferences {
				s := fmt.Sprintf("%s (%s)", owner.Name, owner.Kind)
				o = append(o, s)
			}
			util.AddArray(t, "    Controlled by:", o)
		}
	}
}
