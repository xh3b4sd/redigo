module github.com/xh3b4sd/redigo

go 1.22

require (
	github.com/FZambia/sentinel v1.1.1
	github.com/go-redsync/redsync/v4 v4.13.0
	github.com/gomodule/redigo v1.9.2
	github.com/rafaeljusto/redigomock v2.4.0+incompatible
	github.com/xh3b4sd/breakr v0.1.0
	github.com/xh3b4sd/tracer v0.11.1
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
)

retract [v0.0.0, v0.34.0]
