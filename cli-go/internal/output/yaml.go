package output

import (
	"encoding/json"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(w io.Writer, data interface{}, _ *TableDef) error {
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(j, &doc); err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	clearYAMLStyle(&doc)
	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	_, err = fmt.Fprint(w, string(out))
	return err
}

func clearYAMLStyle(n *yaml.Node) {
	if n.Kind == yaml.ScalarNode && n.Tag == "!!str" {
		n.Encode(n.Value)
	} else {
		n.Style = 0
	}
	for _, c := range n.Content {
		clearYAMLStyle(c)
	}
}
