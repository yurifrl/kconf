package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// T ...
type T struct {
	APIVersion     string      `yaml:"apiVersion"`
	CurrentContext string      `yaml:"current-context"`
	Kind           string      `yaml:"kind"`
	Preferences    interface{} `yaml:"preferences,omitempty"`
	Clusters       []struct {
		Name    string      `yaml:"name"`
		Cluster interface{} `yaml:"cluster"`
	} `yaml:"clusters"`
	Contexts []struct {
		Name    string      `yaml:"name"`
		Context interface{} `yaml:"context"`
	} `yaml:"contexts"`
	Users []struct {
		Name string      `yaml:"name"`
		User interface{} `yaml:"user"`
	} `yaml:"users"`
}

func main() {
	// configs := fmt.Sprintf("/keybase/private/%s/kconf", os.Getenv("USER"))
	configs := fmt.Sprintf("/keybase/private/%s/kconf", "yurifrl")
	config := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
	files, err := ioutil.ReadDir(configs)
	if err != nil {
		log.Fatal(err)
	}

	master := T{
		APIVersion:     "v1",
		CurrentContext: "",
		Kind:           "Config",
	}
	buffer := T{}
	for _, f := range files {
		file := fmt.Sprintf("%s/%s", configs, f.Name())
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(bs, &buffer); err != nil {
			panic(err)
		}

		// Current context will be the last
		master.CurrentContext = buffer.CurrentContext
		master.Clusters = append(master.Clusters, buffer.Clusters[0])
		master.Contexts = append(master.Contexts, buffer.Contexts[0])
		master.Users = append(master.Users, buffer.Users[0])
	}

	bs, err := yaml.Marshal(master)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(config, bs, 0644); err != nil {
		panic(err)
	}
	fmt.Println("File merged")
}
