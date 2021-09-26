package telnet

type NvtCommand struct {
	Command []byte
}

func NewNvt(b []byte) *NvtCommand {
	return &NvtCommand{Command: b}
}

func (n *NvtCommand) WILL() []byte {
	n.Command[1] = 251
	return n.Command
}

func (n *NvtCommand) WONT() []byte {
	n.Command[1] = 252
	return n.Command
}

func (n *NvtCommand) DO() []byte {
	n.Command[1] = 253
	return n.Command
}

func (n *NvtCommand) DONT() []byte {
	n.Command[1] = 254
	return n.Command
}
