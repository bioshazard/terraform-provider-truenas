// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// replaceIfStorePrivatekeyChanged requires replacement only when the EFFECTIVE
// store_privatekey value changes. An unset (null) value means the documented
// default, true, so upgrading from a state written before store_privatekey
// existed (null -> true, once the Default fills it) is NOT a change and must
// not replace the certificate. An unconditional RequiresReplace() did exactly
// that — it read null -> true as a change and planned a destroy/recreate of
// every certificate on upgrade (and, for a CSR, the ACME certificate chained
// to it). A real flip — true<->false, or null->false — still requires
// replacement, since the key cannot be added to or removed from an existing
// certificate in place. Mirrors the replaceIfChangedFromKnown pattern in
// zvol/dataset. (#49 regression; reported in PR #59)
func replaceIfStorePrivatekeyChanged() planmodifier.Bool {
	return boolplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.BoolRequest, resp *boolplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = storesPrivatekey(req.StateValue) != storesPrivatekey(req.PlanValue)
		},
		"Replaced only when the effective store_privatekey value changes (unset reads as the default true).",
		"Replaced only when the effective store_privatekey value changes (unset reads as the default true).",
	)
}
