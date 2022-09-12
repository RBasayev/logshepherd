package main

type inputDef struct {
	ID         string
	Input      string
	Output     string
	DumpBuffer int      `yaml:"dump_buffer"`
	DumpUpon   []string `yaml:"dump_upon,omitempty"`
	FullOutput bool     `yaml:"full_output,omitempty"`
	RotateAt   int      `yaml:"rotate_at,omitempty"`
	Filters    map[string][]string
}

type logshepherdConf struct {
	Reload     string
	Threads    int
	OutputFull map[string]string `yaml:"output_full,omitempty"`
	Routes     []inputDef
}
