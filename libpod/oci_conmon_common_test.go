//go:build !remote && (linux || freebsd)

package libpod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConmonOCIRuntimeVersionInfo(t *testing.T) {
	writeCommand := func(name, output string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), name)
		script := "#!/bin/sh\nprintf '%b' '" + output + "'\n"
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
		return path
	}

	runtime := &ConmonOCIRuntime{
		name:       "test-runtime",
		path:       writeCommand("runtime", "runtime version 1.2\\n"),
		conmonPath: writeCommand("conmon", "conmon version 2.1\\nbuild id 3\\n"),
	}

	conmonVersion, runtimeVersion, err := runtime.versionInfo()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := conmonVersion, "conmon version 2.1, build id 3"; got != want {
		t.Errorf("conmon version = %q, want %q", got, want)
	}
	if got, want := runtimeVersion, "runtime version 1.2"; got != want {
		t.Errorf("runtime version = %q, want %q", got, want)
	}

	info, err := (&Runtime{defaultOCIRuntime: runtime}).VersionInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info.Kernel == "" || info.Distribution.Distribution == "" {
		t.Errorf("VersionInfo omitted host data: %+v", info)
	}
	if info.ConmonVersion != conmonVersion || info.OCIRuntimeVersion != runtimeVersion {
		t.Errorf("VersionInfo returned wrong executable versions: %+v", info)
	}
}
