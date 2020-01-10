package images

import (
	"bytes"
	"drone-release/tools/utils"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
)

// SyncRecordFile is var
var SyncRecordFile string

const syncConfigString = `{
  "auth": {
    "{{.SrcRegistry}}": {
      "username": "admin",
      "password": "nsf0cus.@k8s",
      "insecure": true
    },
    "{{.RemoteRegistry}}": {
      "username": "admin",
      "password": "nsf0cus.@k8s",
      "insecure": true
    }
  },
  "images": {
  	{{- range $i, $x := .Modules}}
  	"{{$x.SrcRegistry}}{{$x.Prefix}}{{$x.ImageName}}:{{$x.ImageVersion}}": "{{$x.RemoteRegistry}}{{$x.Prefix}}{{$x.ImageName}}:{{$x.ImageVersion}}"{{$x.EndOfLine}}
  	{{- end}}
  }
}`

type module struct {
	SrcRegistry    string
	RemoteRegistry string
	ImageName      string
	Prefix         string
	ImageVersion   string
	EndOfLine      string
}

// GenerateSyncConfig is a function to generate image-sync config
func GenerateSyncConfig() {
	cg := utils.ParseRecordTomlConfig(SyncRecordFile)
	t := template.Must(template.New("tmpl").Parse(syncConfigString))
	envMap := make(map[string]interface{})
	var moduleList []module
	for i := range cg.Modules {
		var curMoudle module
		curMoudle.SrcRegistry = cg.Drone.SrcRegistry
		curMoudle.RemoteRegistry = cg.Drone.RemoteRegistry
		curMoudle.ImageName = cg.Modules[i].ImageName
		curMoudle.Prefix = cg.Modules[i].Prefix
		if len(cg.Modules[i].TemplateDict) > 0 {
			envMap := make(map[string]string)
			for j := range cg.Modules[i].TemplateDict {
				cMap := cg.Modules[i].TemplateDict[j].(map[string]interface{})
				for k, v := range cMap {
					envMap[k] = v.(string)
				}
			}
			curMoudle.ImageVersion = envMap["TEMPLATE_IMAGE_VERSION"]
		} else {
			curMoudle.ImageVersion = cg.Modules[i].SvnVersion
		}
		if i == len(cg.Modules)-1 {
			curMoudle.EndOfLine = ""
		} else {
			curMoudle.EndOfLine = ","
		}
		moduleList = append(moduleList, curMoudle)
	}
	envMap["SrcRegistry"] = cg.Drone.SrcRegistry
	envMap["RemoteRegistry"] = cg.Drone.RemoteRegistry
	envMap["Modules"] = moduleList
	var buffer bytes.Buffer
	_ = t.Execute(&buffer, envMap)
	imageSvncPath := cg.Drone.ImageSvncPath
	cmd := fmt.Sprintf("rm -rf %v && mkdir -p %v", imageSvncPath, imageSvncPath)
	_, cmdErr := utils.ExecShell(cmd)
	if "" != cmdErr {
		log.Panicln(cmdErr)
	}
	imageSyncConfigFilePath := fmt.Sprintf("%v/%v", imageSvncPath, "config.json")
	_ = ioutil.WriteFile(imageSyncConfigFilePath, buffer.Bytes(), 0755)

}
