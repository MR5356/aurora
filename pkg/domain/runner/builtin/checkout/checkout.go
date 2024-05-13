package main

import (
	"github.com/MR5356/aurora/pkg/domain/runner/proto"
	"github.com/MR5356/aurora/pkg/domain/runner/shared"
)

type CheckoutTask struct {
	shared.UnimplementedITask
}

func (t *CheckoutTask) GetInfo() *proto.TaskInfo {
	return &proto.TaskInfo{
		Label:       "Checkout",
		Abstract:    "Checkout a Git repository at a particular version",
		Author:      "Rui Ma",
		DownloadUrl: "",
		ProjectUrl:  "",
		Icon:        "",
		Version:     "v1.0.0",
		Usage:       "",
	}
}

func (t *CheckoutTask) GetParams() *proto.TaskParams {
	return &proto.TaskParams{
		Params: []*proto.TaskParam{
			{
				Title:       "Repository",
				Placeholder: "repository address",
				Order:       1,
				Type:        "string",
				Required:    true,
				Key:         "repository",
				Value:       "",
			},
			{
				Title:       "Branch",
				Placeholder: "repository branch",
				Order:       2,
				Type:        "string",
				Required:    true,
				Key:         "branch",
				Value:       "",
			},
			{
				Title:       "Submodules",
				Placeholder: "whether to download submodules",
				Order:       3,
				Type:        "switch",
				Required:    false,
				Key:         "submodules",
				Value:       "false",
			},
			{
				Title:       "Token",
				Placeholder: "token",
				Order:       4,
				Type:        "string",
				Required:    false,
				Key:         "token",
				Value:       "",
			},
		},
	}
}

func (t *CheckoutTask) SetParams(params *proto.TaskParams) {}

func (t *CheckoutTask) Start() error {
	return nil
}

func (t *CheckoutTask) Stop() error {
	return nil
}

func (t *CheckoutTask) Pause() error {
	return nil
}

func (t *CheckoutTask) Resume() error {
	return nil
}

func main() {
	shared.Serve(&CheckoutTask{})
}
