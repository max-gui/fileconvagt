package filesets

import (
	"context"
	"flag"

	"github.com/max-gui/logagent/pkg/confload"
)

var (
	Nexusurl, Nexususr, Nexuspass *string
	// Cacheminutes                   *int
)

func StartupInit(nexusurl, nexususr, nexuspass string) {

	*Nexusurl = nexusurl
	*Nexususr = nexususr
	*Nexuspass = nexuspass
	// confmap := map[string]interface{}{}
	// yaml.Unmarshal(bytes, confmap)
	// *Acltoken = acltoken // confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["private"].(map[interface{}]interface{})["acl-token"].(string)
}

func StartupInitWithConfMap(confmap map[string]interface{}, c context.Context) {
	nexusConf := confload.GetConfDetail("nexus", confmap, c)
	if nexusConf != nil {
		nexusconfMap := nexusConf.(map[string]interface{})
		nexusurl := nexusconfMap["url"].(string)
		nexususr := nexusconfMap["username"].(string)
		nexuspass := nexusconfMap["password"].(string)
		StartupInit(nexusurl, nexususr, nexuspass)
	}
}

func init() {
	// Appname = "charon"
	Nexusurl = flag.String("nexus url", "", "nexus url")
	Nexususr = flag.String("nexus user", "", "nexus user")
	Nexuspass = flag.String("nexus pass", "", "nexus pass")
}

// var Reppath = func() string {
// 	return *Apppath + PthSep + *Repopathname + PthSep
// }
