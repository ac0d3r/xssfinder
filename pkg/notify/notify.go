package notify

import (
	"io"
	"os"

	"github.com/Buzz2d0/xssfinder/pkg/notify/dingbot"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

type Notifier interface {
	Notify(title, text string) error
}

type Notifiers []Notifier

var _ Notifier = Notifiers{}

func (n Notifiers) Notify(title, text string) error {
	var eg errgroup.Group
	for _, notifier := range n {
		nn := notifier
		eg.Go(func() error {
			return nn.Notify(title, text)
		})
	}
	return eg.Wait()
}

type notifiersConfig struct {
	Dingbot dingbot.Dingbot `yaml:"dingbot" json:"dingbot"`
}

func NewNotifierWithYaml(file string) (Notifier, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	t := &notifiersConfig{}
	if err := yaml.Unmarshal(data, t); err != nil {
		return nil, err
	}
	return Notifiers{dingbot.New(t.Dingbot.Token, t.Dingbot.Secret)}, nil
}
