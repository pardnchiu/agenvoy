package discordCommand

// * use static types for ensure type safety and avoid typos
type CommandType int

const (
	CmdHelp CommandType = iota
	CmdRole
)

var commands = []CommandType{
	CmdHelp,
	CmdRole,
}

func (c CommandType) Text() string {
	switch c {
	case CmdHelp:
		return "help"
	case CmdRole:
		return "role"
	default:
		return ""
	}
}

func getCmd(cmd string) CommandType {
	switch cmd {
	case "help", "/help":
		return CmdHelp
	case "role", "/role":
		return CmdRole
	default:
		return -1
	}
}
