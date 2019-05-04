package runner

// InputBuilder represents a builder for the inputs to a Runner.Run.
type InputBuilder struct {
	program string
	dir     string
	args    []string
	env     map[string]string
}

// Creates an InputBuilder for the given command.
func NewInputBuilder(program string) *InputBuilder {
	return &InputBuilder{program: program}
}

// Change what directory the command should be run in.
func (b *InputBuilder) WithDir(dir string) *InputBuilder {
	b.dir = dir
	return b
}

// Include what arguments to pass to the command.
func (b *InputBuilder) WithArgs(args ...string) *InputBuilder {
	b.args = args
	return b
}

// Include what the environment should be when the command is run.
func (b *InputBuilder) WithEnv(env map[string]string) *InputBuilder {
	b.env = env
	return b
}

// Return the inputs in a format that a Runner.Run can consume.
func (b *InputBuilder) Build() (string, []string, string, map[string]string) {
	return b.program, b.args, b.dir, b.env
}
