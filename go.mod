module github.com/satyshef/checker

go 1.18

//replace github.com/satyshef/tdlib => ../tdlib

//replace github.com/satyshef/tdbot => ../tdbot

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/satyshef/mslib/unimes v0.0.0-20220525151958-4576b99c5016
	github.com/satyshef/tdbot v0.2.37
	github.com/satyshef/tdlib v0.2.25
)

require (
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)
