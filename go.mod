module github.com/max-gui/fileconvagt

go 1.15

require (
	github.com/max-gui/logagent v0.0.0-00010101000000-000000000000
	github.com/max-gui/redisagent v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.5.1
	gopkg.in/yaml.v2 v2.2.8

)

replace github.com/max-gui/logagent => ../logagent

replace github.com/max-gui/redisagent => ../redisagent
