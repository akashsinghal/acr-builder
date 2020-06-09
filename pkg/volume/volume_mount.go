// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package volume

import (
	"github.com/pkg/errors"
)

//ValueMount describes a single value to be mounted as a file
type ValueMount struct {
	FileName string `yaml:"filename"`
	Value    string `yaml:"value"`
	B64dec   bool   `yaml:"b64dec"`
}

// VolumeMount describes a Docker bind mounted volume.
type VolumeMount struct {
	Name   string        `yaml:"name"`
	Values []*ValueMount `yaml:"values"`
}

//Validate checks whether VolumeMount is well formed
func (v *VolumeMount) Validate() error {
	if v == nil {
		return nil
	}
	if v.Name == "" || len(v.Values) <= 0 {
		return errors.New("volume name or values is empty")
	}
	for _, valueMount := range v.Values {
		if valueMount != nil {
			if valueMount.FileName == "" {
				return errors.New("filename provided for value is empty")
			}
		}
	}
	return nil
}
