package gvk

// GVK represents the Group-Version-Kind structure for extracting
// apiVersion and kind from manifest YAML files.
// Note: This struct intentionally only includes the fields needed for
// kind extraction. Including Metadata would cause YAML unmarshal failures
// when manifests contain relationships with protobuf enum types.
type GVK struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}
