package git

type IGitHandler interface {
	GetCommitMessages(projectPath string) []string
}
