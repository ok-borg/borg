package conf

import (
	flag "github.com/juju/gnuflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
)

var (
	F = flag.Bool("f", false, "Print full results, ie. no more '...'")
	L = flag.Int("l", 5, "Result list limit. Defaults to 5")
	H = flag.String("h", "borg.crufter.com", "Server to connect to")
	P = flag.Bool("p", false, "Private search. Your search won't leave a trace. Pinky promise. Don't use this all the time if you want to see the search result relevancy improved")
)

var (
	HomeDir string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	HomeDir = usr.HomeDir
	os.Mkdir(HomeDir+"/.borg", os.ModePerm)
	os.Create(HomeDir + "/.borg/edit")
}

type Config struct {
	Token       string
	DefaultTags []string
}

func (c Config) Save() error {
	bs, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(HomeDir+"/.borg/config.yml", bs, os.ModePerm)
}

func Get() (Config, error) {
	bs, err := ioutil.ReadFile(HomeDir + "/.borg/config.yml")
	if err != nil {
		panic(err)
	}
	c := Config{}
	return c, yaml.Unmarshal(bs, &c)
}
