package template

import (
	"bytes"
	"drone-release/tools/utils"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

var waitgroup sync.WaitGroup

// TemplateRecordFile is var
var TemplateRecordFile string

func templateDeployYaml(cg utils.RecordTomlConfig) {
	codePath := cg.Drone.CodePath
	deployDirPath := cg.Drone.DeployDirPath
	cmd := fmt.Sprintf("rm -rf %v && mkdir -p %v", deployDirPath, deployDirPath)
	_, cmdErr := utils.ExecShell(cmd)
	if "" != cmdErr {
		log.Panicln(cmdErr)
	}
	for item := range cg.Modules {
		recordInfo := cg.Modules[item]
		waitgroup.Add(1)
		go func() {
			for item := range recordInfo.DeployFilePath {
				filePath := fmt.Sprintf("%v/%v/%v", codePath, recordInfo.ModuleName, recordInfo.DeployFilePath[item])
				// fileContent, err := ioutil.ReadFile(filePath)
				// if !strings.Contains(string(fileContent), "{{") || !strings.Contains(string(fileContent), "}}") {
				// 	//not need replace from template
				// 	continue
				// }
                
                //只渲染后缀为.tmpl的文件， 其余文件直接cp
                if filePath[len(filePath)-5:] != ".tmpl" {
                    log.Printf("file_name: %v", filePath)
                    cmd := fmt.Sprintf("cp %v %v", filePath, deployDirPath)
		            _, cmdErr := utils.ExecShell(cmd)
		            if "" != cmdErr {
		            	log.Panicln(cmdErr)
		            }
                    continue
                }

				tmpl, err := template.ParseFiles(filePath)
				if err != nil {
					log.Printf("template parse failed:%v", err)
					continue
				}
				var buffer bytes.Buffer
				splits := strings.Split(recordInfo.DeployFilePath[item], "/")
				deployFilePath := fmt.Sprintf("%v/%v", deployDirPath, strings.Replace(splits[len(splits)-1], ".tmpl", "", -1))
				log.Printf("deployFilePath %v", deployFilePath)
				envMap := make(map[string]string)
				if len(recordInfo.TemplateDict) > 0 {
					for i := range recordInfo.TemplateDict {
						cMap := recordInfo.TemplateDict[i].(map[string]interface{})
						for k, v := range cMap {
							envMap[k] = v.(string)
						}
					}
				} else {
					envMap["TEMPLATE_IMAGE_VERSION"] = recordInfo.SvnVersion
				}
				_ = tmpl.Execute(&buffer, envMap)
                t := buffer.Bytes()
                if cg.Drone.Name == "idc" {
                    m := strings.Replace(string(t), "annotations: {kubernetes.io/ingress.class: kong-nuri}", "annotations: {kubernetes.io/ingress.class: kong-tembin}", -1)
				    _ = ioutil.WriteFile(deployFilePath, []byte(m), 0755)
               } else {
				    _ = ioutil.WriteFile(deployFilePath, t, 0755)
               }
			}
			waitgroup.Done()
		}()
	}
}

//Main is config command
func Main() {
	cg := utils.ParseRecordTomlConfig(TemplateRecordFile)
	templateDeployYaml(cg)
	waitgroup.Wait()
}
