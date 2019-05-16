package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
	yaml "gopkg.in/yaml.v2"
)

// KubeConfig ...
type kubeConfig struct {
	APIVersion     string      `yaml:"apiVersion"`
	CurrentContext string      `yaml:"current-context"`
	Kind           string      `yaml:"kind"`
	Preferences    interface{} `yaml:"preferences,omitempty"`
	Clusters       []struct {
		Name    string      `yaml:"name"`
		Cluster interface{} `yaml:"cluster"`
	} `yaml:"clusters,omitempty"`
	Contexts []struct {
		Name    string      `yaml:"name"`
		Context interface{} `yaml:"context"`
	} `yaml:"contexts,omitempty"`
	Users []struct {
		Name string      `yaml:"name"`
		User interface{} `yaml:"user"`
	} `yaml:"users,omitempty"`
}

var (
	version = "1.0.1"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
			Value:   filepath.FromSlash(usr.HomeDir + "/.kconf/config.yaml"),
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "configs",
			Usage:   "Configs source file for kubebernetes configs",
			EnvVars: []string{"CONFIGS"},
			Value:   filepath.FromSlash(usr.HomeDir + "/.kconf/configs"),
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "kubernetes.config",
			Usage:   "Path for kubernetes config file that will receive the sources",
			EnvVars: []string{"KUBE_CONFIG"},
			Value:   filepath.FromSlash(usr.HomeDir + "/.kube/config"),
		}),
	}
	app := &cli.App{
		Action:  run,
		Version: version,
		Before:  altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:   flags,
	}

	app.Run(os.Args)
}

func run(c *cli.Context) (err error) {
	configs := c.String("configs")
	config := c.String("kubernetes.config")
	files, err := ioutil.ReadDir(configs)
	if err != nil {
		log.Fatal(err)
	}

	master := kubeConfig{
		APIVersion:     "v1",
		CurrentContext: "",
		Kind:           "Config",
	}
	buffer := kubeConfig{}
	for _, f := range files {
		file := fmt.Sprintf("%s/%s", configs, f.Name())
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(bs, &buffer); err != nil {
			return err
		}

		// Current context will be the last
		if len(buffer.CurrentContext) > 0 {
			master.CurrentContext = buffer.CurrentContext
		}
		if buffer.Preferences != "" {
			master.Preferences = buffer.Preferences
		}
		if len(buffer.Clusters) > 0 {
			master.Clusters = append(master.Clusters, buffer.Clusters[0])
		}
		if len(buffer.Contexts) > 0 {
			master.Contexts = append(master.Contexts, buffer.Contexts[0])
		}
		if len(buffer.Users) > 0 {
			master.Users = append(master.Users, buffer.Users[0])
		}
	}

	if len(master.Contexts) > 0 {
		master.CurrentContext = master.Contexts[0].Name
	}
	bs, err := yaml.Marshal(master)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(config, bs, 0644); err != nil {
		return err
	}
	log.Println("File merged")
	return err
}
