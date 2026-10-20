package usertypes

type WtfigType struct {
	Version   string `yaml:"version"`
	ApiSwitch bool   `yaml:"apiSwitch"`
	// ApiDetail []Apifig `yaml:"apiDetail"`
	ApiDetail    []Apifig `yaml:"apiDetail"`
	Blocks       Blockfig `yaml:"blocks"`
	ContentBytes []byte   `yaml:"-" json:"-"`
	Service      string   `yaml:"service"`
	From         string   `yaml:"from"`
	WtPath       string   `yaml:"wtpath"`
}

// type WtfigType struct {
// 	Version   string `yaml:"version"`
// 	ApiSwitch bool   `yaml:"apiSwitch"`
// 	// ApiDetail []Apifig `yaml:"apiDetail"`
// 	ApiDetail    map[string]Apifig `yaml:"apiDetail"`
// 	Blocks       Blockfig          `yaml:"blocks"`
// 	ContentBytes []byte            `yaml:"-" json:"-"`
// 	Service      string            `yaml:"service"`
// }

type Apifig struct {
	Id      string     `yaml:"id"`
	Path    string     `yaml:"path"`
	Pattern string     `yaml:"pattern"`
	Zone    Zonefig    `yaml:"zone"`
	Service Servicefig `yaml:"service"`
	Limit   Limitfig   `yaml:"limit"`
}

type Zonefig struct {
	Block bool     `yaml:"block"`
	White []string `yaml:"white"`
}

type Servicefig struct {
	Block bool     `yaml:"block"`
	White []string `yaml:"white"`
}

type Limitfig struct {
	Block bool `yaml:"block"`
	Qps   int  `yaml:"qps"`
	Burst int  `yaml:"burst"`
}

type Blockfig struct {
	Users     []string `yaml:"users"`
	Source_ip []string `yaml:"source_ip"`
}
