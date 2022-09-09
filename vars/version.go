package vars

import "fmt"

type Version struct {
	major int
	minor int
	patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func Current() *Version {
	return &Version{major: 0, minor: 8, patch: 0}
}
