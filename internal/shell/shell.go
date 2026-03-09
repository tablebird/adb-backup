package shell

type Shell interface {
	RunCommand(cmd string, args ...string) (string, error)
}
