package prompt

type PromptSpec struct {
	Role            string   `json:"role,omitempty" yaml:"role,omitempty"`
	Objective       string   `json:"objective,omitempty" yaml:"objective,omitempty"`
	Context         string   `json:"context,omitempty" yaml:"context,omitempty"`
	Scope           []string `json:"scope,omitempty" yaml:"scope,omitempty"`
	Constraints     []string `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	OutputSpec      []string `json:"output_spec,omitempty" yaml:"output_spec,omitempty"`
	QualityCriteria []string `json:"quality_criteria,omitempty" yaml:"quality_criteria,omitempty"`
}
