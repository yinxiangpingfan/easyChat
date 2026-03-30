package config

import (
	"path"
	"testing"
)

func TestGetConfig(t *testing.T) {
	config1 := GetConfig(1, path.Join("..", "local.env"))
	config := GetConfig(0, "")
	t.Logf("%v", config1)
	t.Logf("%v", config)
}
