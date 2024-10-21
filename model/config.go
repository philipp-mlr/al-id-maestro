package model

type Config struct {
	RemoteConfiguration []RemoteConfiguration `yaml:"repositories"`
	ConfigIDRanges      []ConfigIDRange       `yaml:"idRanges"`
}

type RemoteConfiguration struct {
	RepositoryName  string   `yaml:"name"`
	RepositoryURL   string   `yaml:"url"`
	RemoteName      string   `yaml:"remoteName"`
	GithubAuthToken string   `yaml:"authToken"`
	ExcludeBranches []string `yaml:"excludeBranches"`
}

type ConfigIDRange struct {
	ObjectType ObjectType `yaml:"objectType"`
	StartID    uint       `yaml:"from"`
	EndID      uint       `yaml:"to"`
}
