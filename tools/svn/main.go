package svn

import (
	"drone-release/tools/utils"
	"fmt"
	"log"
	"sync"
)

// SvnUserName is var
var SvnUserName string

// SvnPassword is var
var SvnPassword string

// SvnRecordFile is var
var SvnRecordFile string

var waitgroup sync.WaitGroup

func svnCheckout(cg utils.RecordTomlConfig) {
	f0 := func(codePath string) {
		cmd := fmt.Sprintf("rm -rf %v && mkdir -p %v", codePath, codePath)
		log.Printf(cmd)
		_, cmdErr := utils.ExecShell(cmd)
		if "" != cmdErr {
			log.Panicln(cmdErr)
		}
	}
	f1 := func(moduleName string, svnPath string, svnVersion string, svnLocalPath string) {
		log.Printf("svn checkout module_name:%v, svn_path: %v, svn_version: %v", moduleName, svnPath, svnVersion)
		cmd := fmt.Sprintf(
			"svn checkout %v %v/%v --revision %v --username %v --password %v --no-auth-cache --non-interactive --trust-server-cert-failures=unknown-ca,cn-mismatch,expired,not-yet-valid,other",
			svnPath, svnLocalPath, moduleName, svnVersion, SvnUserName, SvnPassword)
		// log.Printf(cmd)
		_, cmdErr := utils.ExecShell(cmd)
		if "" != cmdErr {
			log.Panicln(cmdErr)
		}
		waitgroup.Done()
	}
	codePath := cg.Drone.CodePath
	f0(codePath)
	for item := range cg.Modules {
		recordInfo := cg.Modules[item]
		go f1(recordInfo.ModuleName, recordInfo.SvnPath, recordInfo.SvnVersion, codePath)
		waitgroup.Add(1)
	}
}

// Main svn checkout function
func Main() {
	cg := utils.ParseRecordTomlConfig(SvnRecordFile)
	svnCheckout(cg)
	waitgroup.Wait()
}
