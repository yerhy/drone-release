package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/toml"
)

// RecordTomlConfig is a golang struct of svn_record
type RecordTomlConfig struct {
	Drone struct {
        Name           string `toml:name`
		CodePath       string `toml:"code_path"`
		DeployDirPath  string `toml:"deploy_dir_path"`
		ImageSvncPath  string `toml:"image_svnc_path"`
		SrcRegistry    string `toml:"src_registry"`
		RemoteRegistry string `toml:"remote_registry"`
	} `toml:"drone"`
	Modules []struct {
		ModuleName     string        `toml:"module_name"`
		ImageName      string        `toml:"image_name"`
		Prefix         string        `toml:"prefix"`
		SvnPath        string        `toml:"svn_path"`
		SvnVersion     string        `toml:"svn_version"`
		DeployFilePath []string      `toml:"deploy_file_path"`
		TemplateDict   []interface{} `toml:"template_dict"`
	} `toml:"modules"`
}

// ExecShell : shell executor
func ExecShell(s string) (string, string) {
	cmd := exec.Command("/bin/sh", "-c", s)
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	err := cmd.Run()
	if err != nil {
		log.Panicln(fmt.Sprintf("exec cmd '%v' error:%v", s, cmdErr.String()))
	}
	return cmdOut.String(), cmdErr.String()
}

// ParseRecordTomlConfig is a function to transform toml to golang struct
func ParseRecordTomlConfig(tomlConfigPath string) RecordTomlConfig {
	var cg RecordTomlConfig
	if _, err := toml.DecodeFile(tomlConfigPath, &cg); err != nil {
		log.Fatal(err)
	}
	return cg
}
