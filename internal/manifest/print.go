package manifest

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/protobufyaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Print(input proto.Message) error {
	marshalJsonBytes, err := protojson.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "failed to yaml marshalJsonBytes")
	}
	marshalYamlBytes, err := protobufyaml.JSONToYAML(marshalJsonBytes)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json to yaml")
	}

	fmt.Printf("%v\n", string(marshalYamlBytes))
	return nil
}
