package nidhogg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeConfig(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	assert.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestGetConfigReadsYAMLAndJSON(t *testing.T) {
	yamlConf := `daemonsets:
  - name: kiam
    namespace: kube-system
nodeSelector:
  - "node-role.kubernetes.io/node"
  - "!node-role.kubernetes.io/master"
taintNamePrefix: "nidhogg.uswitch.com"
taintEffect: NoExecute
taintRemovalDelayInSeconds: 10
`
	jsonConf := `{"daemonsets":[{"name":"kiam","namespace":"kube-system"}],"nodeSelector":["node-role.kubernetes.io/node","!node-role.kubernetes.io/master"],"taintNamePrefix":"nidhogg.uswitch.com","taintEffect":"NoExecute","taintRemovalDelayInSeconds":10}`

	for name, content := range map[string]string{"config.yaml": yamlConf, "config.json": jsonConf} {
		cfg, err := GetConfig(writeConfig(t, name, content))
		assert.NoError(t, err, name)
		assert.Equal(t, []Daemonset{{Name: "kiam", Namespace: "kube-system"}}, cfg.Daemonsets, name)
		assert.Equal(t, []string{"node-role.kubernetes.io/node", "!node-role.kubernetes.io/master"}, cfg.NodeSelector, name)
		assert.Equal(t, "nidhogg.uswitch.com", cfg.TaintNamePrefix, name)
		assert.Equal(t, "NoExecute", cfg.TaintEffect, name)
		assert.Equal(t, 10, cfg.TaintRemovalDelayInSeconds, name)
		assert.Len(t, cfg.DaemonsetSelectors, 1, name)
	}
}

func TestGetConfigRejectsEmptyDaemonsets(t *testing.T) {
	_, err := GetConfig(writeConfig(t, "config.yaml", "taintEffect: NoSchedule\n"))
	assert.Error(t, err)
}

func TestGetConfigRejectsUnknownTaintEffect(t *testing.T) {
	_, err := GetConfig(writeConfig(t, "config.yaml", "daemonsets:\n  - name: a\n    namespace: b\ntaintEffect: Nope\n"))
	assert.Error(t, err)
}
