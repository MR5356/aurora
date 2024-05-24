package pluginhub

import "github.com/MR5356/aurora/pkg/domain/runner/proto"

type TaskConfig struct {
	Label       string             `yaml:"label"`
	Description string             `yaml:"description"`
	Author      string             `yaml:"author"`
	Icon        string             `yaml:"icon"`
	Version     string             `yaml:"version"`
	Params      []*proto.TaskParam `yaml:"params"`
}
