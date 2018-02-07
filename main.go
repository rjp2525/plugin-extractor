package main

import (
    "fmt"
    "log"
    "archive/zip"
    "io"
    "io/ioutil"
    "os"
    "gopkg.in/yaml.v2"
    "encoding/json"
    "strings"
)

type Config struct {
    Name string `json:"name" yaml:"name"`
    Description string `json:"description,omitempty" yaml:"description,omitempty"`
    Version string `json:"version" yaml:"version"`
    Commands map[string]Command `json:"commands,omitempty" yaml:"commands,omitempty"`
    Permissions map[string]Permission `json:"permissions,omitempty" yaml:"permissions,omitempty"`
}

type Command struct {
    Description string `json:"description,omitempty" yaml:"description,omitempty"`
    Aliases string `json:"aliases,omitempty" yaml:"aliases,omitempty"`
    Permission string `json:"permission,omitempty" yaml:"permission,omitempty"`
    PermissionMessage string `json:"permission-message,omitempty" yaml:"permission-message,omitempty"`
    Usage string `json:"usage,omitempty" yaml:"usage,omitempty"`
}

type Permission struct {
    Description string `json:"description,omitempty" yaml:"description,omitempty"`
    Default string `json:"default,omitempty" yaml:"default,omitempty"`
    Children map[string]bool `json:"children,omitempty" yaml:"children,omitempty"`
}

func GetPluginInfo(filename string) (error) {
    return readFromFile(filename)
}

func readFromFile(filename string) (error) {
    var err error
    var file *os.File
    var fi os.FileInfo
    var r *zip.Reader

    if file, err = os.Open(filename); err != nil {
        return err
    }
    defer file.Close()

    if fi, err = file.Stat(); err != nil {
        return err
    }

    if r, err = zip.NewReader(file, fi.Size()); err != nil {
        return err
    }

    return readFromReader(r)
}

func readFromReader(r *zip.Reader) (error) {
    for _, f := range r.File {
        if f.Name == "plugin.yml" {
            rc, err := f.Open()
            if err != nil {
                log.Println(err)
                return err
            }

            by, err := ioutil.ReadAll(rc)
            if err != nil {
                log.Println(err)
                return err
            }

            if err == io.EOF {
                //err = nil
                continue
            }

            var config Config
            err = yaml.Unmarshal(by, &config)
            if err != nil {
                log.Println(err)
                return err
            }

            j, err := json.MarshalIndent(config, "", "    ")
            if err != nil {
                log.Println(err)
                return err
            }

            outfile := strings.ToLower(config.Name) + "-" + config.Version + ".json"
            err = ioutil.WriteFile(outfile, j, 0644)
            if err != nil {
                log.Println(err)
                return err
            }

            rc.Close()
        }
    }

    return nil
}

func main() {
    if len(os.Args) == 1 {
        fmt.Println("A .jar plugin is required to extract data.")
        os.Exit(0)
    }

    pluginData := GetPluginInfo(os.Args[1])

    fmt.Println(pluginData)
}