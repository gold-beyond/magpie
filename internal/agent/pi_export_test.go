package agent

import (
	"path/filepath"
	"testing"

	"github.com/yetone/magpie/internal/edit"
)

// #2373: picking a magpie/ model the acting process can't resolve yet (a
// channel just signed in) wrote settings.json but left models.json stale,
// so Pi started on a model its list doesn't have. The pick now always
// refreshes the provider block, as it does for models that resolve.
func TestPiPickRefreshesStaleModels(t *testing.T) {
	home := syncHome(t)
	dir := filepath.Join(home, ".pi", "agent")
	writeFile(t, filepath.Join(dir, "settings.json"), `{
  "defaultProvider": "magpie",
  "defaultModel": "relay/glm-4.6"
}
`)
	writeFile(t, filepath.Join(dir, "models.json"), `{
  "providers": {
    "magpie": {
      "name": "magpie",
      "models": [
        {"id": "relay/glm-4.6", "name": "GLM", "reasoning": false}
      ]
    }
  }
}
`)
	before := readFile(filepath.Join(dir, "models.json"))
	a := pi(home)
	if err := a.Field("model").Set("magpie/gone-provider-xyz/ghost-model"); err != nil {
		t.Fatal(err)
	}
	if v, _ := edit.GetJSON(filepath.Join(dir, "settings.json"), "defaultModel"); v != "gone-provider-xyz/ghost-model" {
		t.Fatalf("defaultModel = %q\n%s", v, readFile(filepath.Join(dir, "settings.json")))
	}
	if got := readFile(filepath.Join(dir, "models.json")); got == before {
		t.Fatalf("pick left models.json stale:\n%s", got)
	}
}
