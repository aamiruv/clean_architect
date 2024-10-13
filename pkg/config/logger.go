package config

type logger struct {
	level            int
	directory        string
	fileCreationMode int
	remoteURL        string
	console          bool
}

func (l logger) Level() int {
	return l.level
}

func (l logger) Directory() string {
	return l.directory
}

func (l logger) FileCreationMode() int {
	return l.fileCreationMode
}

func (l logger) RemoteURL() string {
	return l.remoteURL
}

func (l logger) Console() bool {
	return l.console
}
