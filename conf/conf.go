package conf

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	flag "github.com/juju/gnuflag"
	"gopkg.in/yaml.v2"
)

var (
	// F flag prints full request
	F = flag.Bool("f", false, "Print full results, ie. no more '...'")

	// L flag limit results to a number
	L = flag.Int("l", 5, "Result list limit. Defaults to 5")

	// H flag specifies the host to connect to
	S = flag.String("s", "ok-b.org", "Server to connect to")

	H = flag.Bool("h", false, "Display help")

	Help = flag.Bool("help", false, "Display help, same as -h")

	// P flag enables private search
	P = flag.Bool("p", false, "Private search. Your search won't leave a trace. Pinky promise. Don't use this all the time if you want to see the search result relevancy improved")

	// D flag enables debug mode
	D = flag.Bool("d", false, "Debug mode")
	// DontPipe
	DontPipe = flag.Bool("dontpipe", false, "Flag for internal use - ignore this")
	// Version flag displays current version
	Version = flag.Bool("version", false, "Print version number")
	// V flag displays current version
	V = flag.Bool("v", false, "Print version number")
)
var (
	// EditFile borg edit file.
	EditFile string
	// ConfigFile borg config file.
	ConfigFile string
	// QueryFile borg query file.
	QueryFile string
)

func init() {
	borgDir := borgDir()

	EditFile = filepath.Join(borgDir, "edit")
	ConfigFile = filepath.Join(borgDir, "config.yml")
	QueryFile = filepath.Join(borgDir, "query")

	os.Mkdir(borgDir, os.ModePerm)
	os.Create(EditFile)
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		os.Create(ConfigFile)
	}
	if _, err := os.Stat(QueryFile); os.IsNotExist(err) {
		os.Create(QueryFile)
	}
}

func borgDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(usr.HomeDir, ".borg")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			dir = filepath.Join(xdgConfigHome, "borg")
		} else {
			dir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}

	return dir
}

// Config file
type Config struct {
	Token       string
	DefaultTags []string
	Editor      string
	PipeTo      string
}

// Save config
func (c Config) Save() error {
	bs, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(QueryFile, bs, os.ModePerm)
}

// Get config
func Get() (Config, error) {
	bs, err := ioutil.ReadFile(QueryFile)
	if err != nil {
		panic(err)
	}
	c := &Config{}
	err = yaml.Unmarshal(bs, c)
	if err != nil {
		return *c, err
	}
	if len(c.Editor) == 0 {
		c.Editor = "vim"
	}
	return *c, nil
}
