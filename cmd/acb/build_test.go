// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"testing"

	"github.com/Azure/acr-builder/graph"
	"github.com/Azure/acr-builder/templating"
	"github.com/Azure/acr-builder/util"
)

func TestCreateBuildTask(t *testing.T) {
	buildCmd := &buildCmd{
		dockerfile:  "HelloWorld/Dockerfile",
		credentials: []string{"foo.azurecr.io;user;pw"},
		context:     "src",
		tags:        []string{"foo:latest", "bar/qux"},
		pull:        true,
		noCache:     false,
		dryRun:      true,
		opts: &templating.BaseRenderOptions{
			Registry: "foo.azurecr.io",
		},
	}

	task, err := buildCmd.createBuildTask()
	if err != nil {
		t.Fatalf("failed to create build task, err: %v", err)
	}
	if len(task.Credentials) == 0 {
		t.Fatalf("Expected to create credentials but no credentials were created")
	}

	taskCreds := task.Credentials[0]
	expectedCreds, _ := graph.CreateCredentialFromString(buildCmd.credentials[0])
	if taskCreds.RegistryName != expectedCreds.RegistryName ||
		taskCreds.RegistryUsername != expectedCreds.RegistryUsername ||
		taskCreds.RegistryPassword != expectedCreds.RegistryPassword {
		t.Fatalf("expected %v Creds, got %v", expectedCreds, taskCreds)
	}

	numSteps := len(task.Steps)
	expectedSteps := 1
	if numSteps != expectedSteps {
		t.Fatalf("expected %d steps, got %d", expectedSteps, numSteps)
	}

	// When registry information is provided, the resulting tags will be
	// prefixed with the fully qualified registry's name.
	buildStep := task.Steps[0]
	expectedTags := []string{"foo.azurecr.io/foo:latest", "foo.azurecr.io/bar/qux"}
	if !util.StringSequenceEquals(buildStep.Tags, expectedTags) {
		t.Fatalf("expected %v to be the task's tags but got %v", expectedTags, buildStep.Tags)
	}
	expectedCmd := "--pull -f HelloWorld/Dockerfile -t foo.azurecr.io/foo:latest -t foo.azurecr.io/bar/qux src"
	if expectedCmd != buildStep.Build {
		t.Fatalf("expected %s as the build command, but got %s", expectedCmd, buildStep.Build)
	}
}