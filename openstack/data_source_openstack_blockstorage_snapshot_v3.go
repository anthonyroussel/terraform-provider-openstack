package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBlockStorageSnapshotV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBlockStorageSnapshotV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Computed values
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageSnapshotV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	client, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	listOpts := snapshots.ListOpts{
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		VolumeID: d.Get("volume_id").(string),
	}

	allPages, err := snapshots.List(client, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query openstack_blockstorage_snapshots_v3: %s", err)
	}

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_blockstorage_snapshots_v3: %s", err)
	}

	if len(allSnapshots) < 1 {
		return diag.Errorf("Your openstack_blockstorage_snapshot_v3 query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var snapshot snapshots.Snapshot

	if len(allSnapshots) > 1 {
		recent := d.Get("most_recent").(bool)

		if recent {
			snapshot = dataSourceBlockStorageV3MostRecentSnapshot(allSnapshots)
		} else {
			log.Printf("[DEBUG] Multiple openstack_blockstorage_snapshot_v3 results found: %#v", allSnapshots)

			return diag.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		snapshot = allSnapshots[0]
	}

	dataSourceBlockStorageSnapshotV3Attributes(d, snapshot)
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceBlockStorageSnapshotV3Attributes(d *schema.ResourceData, snapshot snapshots.Snapshot) {
	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("description", snapshot.Description)
	d.Set("size", snapshot.Size)
	d.Set("status", snapshot.Status)
	d.Set("volume_id", snapshot.VolumeID)

	if err := d.Set("metadata", snapshot.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for snapshot %s: %s", snapshot.ID, err)
	}
}
