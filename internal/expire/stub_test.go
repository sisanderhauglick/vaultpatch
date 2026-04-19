package expire_test

import (
	"fmt"

	"github.com/your-org/vaultpatch/internal/vault"
)

type memClient struct {
	data    map[string]map[string]interface{}
	deleted map[string]bool
}

func stubClient(data map[string]map[string]interface{}) *vault.Client {
	m := &memClient{data: data, deleted: map[string]bool{}}
	return vault.NewStubClient(m)
}

func (m *memClient) ReadSecrets(path string) (map[string]interface{}, error) {
	if m.deleted[path] {
		return map[string]interface{}{}, nil
	}
	v, ok := m.data[path]
	if !ok {
		return map[string]interface{}{}, nil
	}
	out := make(map[string]interface{}, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (m *memClient) WriteSecrets(path string, secrets map[string]interface{}) error {
	m.data[path] = secrets
	return nil
}

func (m *memClient) DeleteSecrets(path string) error {
	if _, ok := m.data[path]; !ok {
		return fmt.Errorf("stub: path %q not found", path)
	}
	m.deleted[path] = true
	return nil
}
