package builders

import "path/filepath"

type distributer struct {
	Assignments map[string]func(src, dst string) error
}

func NewDistributer(assignments map[string]func(src, dst string) error) distributer {
	d := distributer{Assignments: map[string]func(src, dst string) error{}}
	for pattern, builder := range assignments {
		d.Assignments[filepath.Clean(pattern)] = builder
	}
	return d
}
