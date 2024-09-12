package main

type commandType string

const (
	CmdSetUsername commandType = "/username"
	CmdJoinRoom    commandType = "/join"
	CmdListRooms   commandType = "/rooms"
	CmdListMembers commandType = "/members"
	CmdSendMessage commandType = "/msg"
	CmdQuit        commandType = "/quit"
)

type command struct {
	commandType commandType
	args        []string
	client      *client
}
