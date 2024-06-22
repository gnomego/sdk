package exec

type CommandArgs struct {
	args []string
	len  int
}

func NewCommandArgs(args ...string) *CommandArgs {
	return &CommandArgs{args: args, len: len(args)}
}

func ParseCommandArgs(args string) *CommandArgs {
	a := SplitArguments(args)
	return &CommandArgs{args: a}
}

func (c *CommandArgs) ToArgs() []string {
	// make a copy of the args
	args := make([]string, len(c.args))
	copy(args, c.args)
	return args
}

func (c *CommandArgs) Parse(args ...string) *CommandArgs {
	c.args = args
	return c
}

func (c *CommandArgs) Clear() *CommandArgs {
	c.args = []string{}
	c.len = 0
	return c
}

func (c *CommandArgs) PrependMap(m map[string]string) *CommandArgs {
	args := []string{}
	for k, v := range m {
		args = append(args, k, v)
	}
	c.args = append(args, c.args...)
	c.len += len(args)
	return c
}

func (c *CommandArgs) AppendMap(m map[string]string) *CommandArgs {
	for k, v := range m {
		c.Append(k, v)
	}
	return c
}

func (c *CommandArgs) Append(args ...string) *CommandArgs {
	c.len += len(args)
	c.args = append(c.args, args...)
	return c
}

func (c *CommandArgs) Len() int {
	return c.len
}

func (c *CommandArgs) Prepend(args ...string) *CommandArgs {
	c.len += len(args)
	c.args = append(args, c.args...)
	return c
}

func (c *CommandArgs) RemoveAt(index int) *CommandArgs {
	if index < 0 || index >= c.len {
		return c
	}
	c.args = append(c.args[:index], c.args[index+1:]...)
	c.len--
	return c
}

func (c *CommandArgs) Set(index int, arg string) *CommandArgs {
	if index < 0 || index >= c.len {
		return c
	}
	c.args[index] = arg
	return c
}

func (c *CommandArgs) Get(index int) string {
	if index < 0 || index >= c.len {
		return ""
	}
	return c.args[index]
}

func (c *CommandArgs) InsertAt(index int, args ...string) *CommandArgs {
	c.args = append(c.args[:index], append(args, c.args[index:]...)...)
	c.len += len(args)
	return c
}
