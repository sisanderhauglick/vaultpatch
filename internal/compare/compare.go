package compare

import (
	"fmt"

	"github.com/vaultpatch/internal/vault"
)

// Result holds the comparison outcome between two Vault paths.
type Result struct {
	SourcePath string
	DestPath   string
	OnlyInSrc  map[string]string
	OnlyInDst  map[string]string
	Differ     map[string][2]string // key -> [src, dst]
	Match      map[string]string
}

// Compare reads secrets from two paths and returns a structured diff.
func Compare(client *vault.Client, srcPath, dstPath string) (*Result, error) {
	if client == nil {
		return nil, fmt.Errorf("compare: client must not be nil")
	}
	if srcPath == "" {
		return nil, fmt.Errorf("compare: source path must not be empty")
	}
	if dstPath == "" {
		return nil, fmt.Errorf("compare: destination path must not be empty")
	}

	src, err := vault.ReadSecrets(client, srcPath)
	if err != nil {
		return nil, fmt.Errorf("compare: read source: %w", err)
	}
	dst, err := vault.ReadSecrets(client, dstPath)
	if err != nil {
		return nil, fmt.Errorf("compare: read destination: %w", err)
	}

	res := &Result{
		SourcePath: srcPath,
		DestPath:   dstPath,
		OnlyInSrc:  make(map[string]string),
		OnlyInDst:  make(map[string]string),
		Differ:     make(map[string][2]string),
		Match:      make(map[string]string),
	}

	for k, sv := range src {
		if dv, ok := dst[k]; ok {
			if sv == dv {
				res.Match[k] = sv
			} else {
				res.Differ[k] = [2]string{sv, dv}
			}
		} else {
			res.OnlyInSrc[k] = sv
		}
	}
	for k, dv := range dst {
		if _, ok := src[k]; !ok {
			res.OnlyInDst[k] = dv
		}
	}
	return res, nil
}
