package integration

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const itestImage = "mcp-exec:itest"

// TestDockerSandbox builds the production image and verifies two image-level
// invariants that require a real container (constitution V): the sandbox has no
// network (FR-010) and ships the declared packages (FR-016).
func TestDockerSandbox(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping docker build in -short mode")
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}
	if err := pool.Client.Ping(); err != nil {
		t.Skipf("docker daemon unreachable: %v", err)
	}

	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	// Build via the docker CLI (BuildKit). The go-dockerclient legacy build API
	// is unreliable here; container run/verify below uses dockertest.
	build := exec.Command("docker", "build", "-t", itestImage, ".")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build image: %v\n%s", err, out)
	}

	t.Run("packages present", func(t *testing.T) {
		code, logs := runOneShot(t, pool, "import yaml, jinja2; print('pkgs-ok')", false)
		if code != 0 || !strings.Contains(logs, "pkgs-ok") {
			t.Errorf("expected packages available (exit 0), got exit=%d logs=%q", code, logs)
		}
	})

	t.Run("no network", func(t *testing.T) {
		script := "import socket; socket.setdefaulttimeout(3); socket.create_connection(('1.1.1.1', 53))"
		code, logs := runOneShot(t, pool, script, true)
		if code == 0 {
			t.Errorf("expected network failure with --network none, but connection succeeded; logs=%q", logs)
		}
	})
}

func runOneShot(t *testing.T, pool *dockertest.Pool, script string, networkNone bool) (int, string) {
	t.Helper()
	hostCfg := &docker.HostConfig{AutoRemove: false}
	if networkNone {
		hostCfg.NetworkMode = "none"
	}
	c, err := pool.Client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      itestImage,
			Entrypoint: []string{"python3", "-c", script},
		},
		HostConfig: hostCfg,
	})
	if err != nil {
		t.Fatalf("create container: %v", err)
	}
	defer func() {
		_ = pool.Client.RemoveContainer(docker.RemoveContainerOptions{ID: c.ID, Force: true})
	}()

	if err := pool.Client.StartContainer(c.ID, nil); err != nil {
		t.Fatalf("start container: %v", err)
	}
	code, err := pool.Client.WaitContainer(c.ID)
	if err != nil {
		t.Fatalf("wait container: %v", err)
	}

	var buf bytes.Buffer
	if err := pool.Client.Logs(docker.LogsOptions{
		Container:    c.ID,
		OutputStream: &buf,
		ErrorStream:  &buf,
		Stdout:       true,
		Stderr:       true,
	}); err != nil {
		t.Fatalf("logs: %v", err)
	}
	return code, buf.String()
}
