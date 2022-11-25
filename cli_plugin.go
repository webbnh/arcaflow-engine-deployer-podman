package podman

import (
	"arcaflow-engine-deployer-podman/internal/cli_wrapper"
	"io"
	"sync"

	"go.arcalot.io/log"
)

type CliPlugin struct {
	wrapper        cli_wrapper.CliWrapper
	lock           *sync.Mutex
	containerImage string
	readIndex      int64
	config         *Config
	logger         log.Logger
	stdin          io.WriteCloser
	stdout         io.ReadCloser
}

// TODO: unwrap the whole config

func (p *CliPlugin) Write(b []byte) (n int, err error) {
	return p.stdin.Write(b)
}

func (p *CliPlugin) Read(b []byte) (n int, err error) {
	return p.stdout.Read(b)
}

func (p *CliPlugin) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.stdout.Close()
	p.stdin.Close()
	return nil
}
