// Package reorder provides functionality to reorder secret keys at one or more
// Vault paths. The desired key order is specified explicitly; any keys not
// listed are preserved at the end in their original relative order.
//
// Reorder supports a dry-run mode that computes the new ordering without
// writing changes back to Vault, making it safe to preview the effect of a
// reorder operation before committing.
package reorder
