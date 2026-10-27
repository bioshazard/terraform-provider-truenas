// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package dataset

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestQuotaUpdateOnlySendsChanges(t *testing.T) {
	for _, tt := range []struct {
		name          string
		before, after types.Int64
		send          bool
	}{
		{"default zero", types.Int64Value(0), types.Int64Value(0), false},
		{"finite unchanged", types.Int64Value(1048576), types.Int64Value(1048576), false},
		{"set finite", types.Int64Value(0), types.Int64Value(1048576), true},
		{"reset to zero", types.Int64Value(1048576), types.Int64Value(0), true},
		{"unset", types.Int64Null(), types.Int64Null(), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			before := DatasetModel{Quota: tt.before, RefQuota: tt.before, Reservation: tt.before}
			after := DatasetModel{Quota: tt.after, RefQuota: tt.after, Reservation: tt.after, Comments: types.StringValue("changed")}
			p := after.changedUpdateAPIPayload(before)
			for _, key := range []string{"quota", "refquota", "reservation"} {
				value, present := p[key]
				if present != tt.send {
					t.Fatalf("%s present=%v, want %v", key, present, tt.send)
				}
				if present && value != tt.after.ValueInt64() {
					t.Fatalf("%s wrong value", key)
				}
			}
			if p["comments"] != "changed" {
				t.Fatal("lost comment update")
			}
		})
	}
}

func TestQuotaUpdateTreatsEachFieldSeparately(t *testing.T) {
	prior := DatasetModel{Quota: types.Int64Value(4), RefQuota: types.Int64Value(2), Reservation: types.Int64Value(1)}
	for _, name := range []string{"quota", "refquota", "reservation"} {
		t.Run(name, func(t *testing.T) {
			planned := prior
			switch name {
			case "quota":
				planned.Quota = types.Int64Value(0)
			case "refquota":
				planned.RefQuota = types.Int64Value(0)
			case "reservation":
				planned.Reservation = types.Int64Value(0)
			}
			payload := planned.changedUpdateAPIPayload(prior)
			for _, key := range []string{"quota", "refquota", "reservation"} {
				value, present := payload[key]
				if key == name {
					if !present || value != int64(0) {
						t.Fatalf("%s: explicit zero missing: %v", key, payload)
					}
				} else if present {
					t.Fatalf("unchanged %s included: %v", key, payload)
				}
			}
		})
	}
}
