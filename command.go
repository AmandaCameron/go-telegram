package telegram

import (
	"log"
	"strings"

	"runtime"

	"fmt"
	"sort"
)

// CommandFlags represents a set of flags that the bot can use to give hints to
// the Help and Method functions.
type CommandFlags uint32

const (
	// CommandFlagNone means no flags whatsoever.
	CommandFlagNone CommandFlags = 1 << iota

	// CommandFlagHidden means this command will not show up in Help's outpur.
	CommandFlagHidden

	// CommandFlagNoGroup means this command will not be accessible from a group
	// messaging context.
	CommandFlagNoGroup
)

// Command represents a command to respond to.
type Command struct {
	Name        string
	Description string

	Flags CommandFlags

	Handle func(Message)
}

// Match takes the given Message and returns true if the message is
// a) Addressing the bot & command
// b) Appropiate to this context
func (cmd *Command) Match(msg Message) bool {
	inp := strings.TrimSpace(msg.Text)
	self := msg.bot.api.Self

	if inp == "/"+cmd.Name {
		if cmd.Flags&CommandFlagNoGroup == CommandFlagNoGroup {
			return !msg.IsGroup()
		}

		return true
	}

	if msg.IsGroup() {
		if cmd.Flags&CommandFlagNoGroup == CommandFlagNoGroup {
			return false
		}

		if inp == "@"+self.UserName+" "+cmd.Name {
			return true
		} else if inp == "/"+cmd.Name+"@"+self.UserName {
			return true
		}
	}

	return false
}

// Commands represents a set of commands
type Commands []Command

// Add adds a new command to the set.
func (cmds *Commands) Add(cmd Command) {
	*cmds = append(*cmds, cmd)
}

// Handle loops through the commands and calls the proper one's Handle method
// if no matching command is found, it will return false, otherwise true.
func (cmds *Commands) Handle(msg Message) bool {
	for _, cmd := range *cmds {
		if cmd.Match(msg) {
			go func() {
				defer func() {
					err := recover()

					if err != nil {
						const size = 64 << 10
						buf := make([]byte, size)
						buf = buf[:runtime.Stack(buf, false)]

						log.Printf("Recovered from crash: %+v\n%s", err, buf)
						msg.ReplyWith("Fatal bot error. Sorry!").Send()
					}
				}()

				cmd.Handle(msg)
			}()

			return true
		}
	}

	return false
}

// Help generates a Help text blob for the given commands.
// pass true to not include commands disabled for groups.
func (cmds *Commands) Help(groupChat bool) string {
	reply := ""

	commands := commandsByName((*cmds)[:])
	sort.Sort(commands)

	for _, cmd := range commands {
		if cmd.Flags&CommandFlagHidden == CommandFlagHidden {
			continue
		}

		if cmd.Flags&CommandFlagNoGroup == CommandFlagNoGroup && groupChat {
			continue
		}

		reply += fmt.Sprintf("/%s - %s\n", cmd.Name, cmd.Description)
	}

	return reply
}

type commandsByName []Command

func (cmds commandsByName) Len() int {
	return len(cmds)
}

func (cmds commandsByName) Less(i, j int) bool {
	return cmds[i].Name < cmds[j].Name
}

func (cmds commandsByName) Swap(i, j int) {
	cmds[i], cmds[j] = cmds[j], cmds[i]
}
