module github.com/0xundef/php-webshell-detect

go 1.18

replace github.com/0xundef/php-webshell-detect => ./

exclude github.com/BurntSushi/toml-test v0.1.1-0.20210723065233-facb9eccd4da

require (
	github.com/deckarep/golang-set/v2 v2.1.0
	github.com/schollz/progressbar/v3 v3.13.1
	github.com/sirupsen/logrus v1.9.0
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/term v0.8.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
