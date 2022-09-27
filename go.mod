module github.com/max-gui/fileconvagt

go 1.18

require (
	github.com/max-gui/logagent v0.0.0-20211102065508-44b5d1757320
	github.com/max-gui/redisagent v0.0.0-20211104054521-b437c64da1c5
	github.com/stretchr/testify v1.7.1
	gopkg.in/yaml.v2 v2.4.0

)

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gomodule/redigo v1.8.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/max-gui/logagent => ../logagent

replace github.com/max-gui/redisagent => ../redisagent
