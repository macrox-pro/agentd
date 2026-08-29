package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

const invalidUserConfigNotifyPrefix = "agentd: invalid user config "

// PrepareUserConfig ensures the user config file exists and compiles before daemon load.
// Creates a minimal bootstrap file when missing. Invalid parse/compile failures print to
// notify and return an error without modifying the file. Only daemon start should call this.
func PrepareUserConfig(userPath string, notify io.Writer) error {
	if userPath == "" {
		return fmt.Errorf("user config path: home directory unavailable")
	}
	st, err := os.Stat(userPath)
	if err != nil {
		if os.IsNotExist(err) {
			return writeBootstrapUserConfig(userPath)
		}
		return fmt.Errorf("stat user config: %w", err)
	}
	if st.IsDir() {
		return fmt.Errorf("read user config: %q is a directory", userPath)
	}
	_, fc, err := readFileConfig(userPath)
	if err != nil {
		if isUserConfigParseError(err) {
			notifyInvalidUserConfig(notify, userPath, err)
			return fmt.Errorf("invalid user config %q: %w", userPath, err)
		}
		return fmt.Errorf("read user config: %w", err)
	}
	if _, _, _, _, err := Compile(fc); err != nil {
		notifyInvalidUserConfig(notify, userPath, err)
		return fmt.Errorf("invalid user config %q: %w", userPath, err)
	}
	return nil
}

func bootstrapUserFileConfig() *fileConfig {
	enabled := true
	return &fileConfig{
		Version: 1,
		Policy: &filePolicy{
			Fail: string(FailClosed),
		},
		Guards: &fileGuards{
			Secrets: &fileSecretsGuard{
				Enabled: &enabled,
				Action:  string(GuardAsk),
			},
		},
	}
}

func writeBootstrapUserConfig(userPath string) error {
	raw, err := yaml.Marshal(bootstrapUserFileConfig())
	if err != nil {
		return fmt.Errorf("marshal user config: %w", err)
	}
	if err := writeYAMLAtomic(userPath, raw); err != nil {
		return fmt.Errorf("write user config: %w", err)
	}
	return nil
}

func notifyInvalidUserConfig(notify io.Writer, userPath string, err error) {
	if notify == nil {
		return
	}
	_, _ = fmt.Fprintf(notify, "%s%s: %v\n", invalidUserConfigNotifyPrefix, userPath, err)
}

func isUserConfigParseError(err error) bool {
	return strings.Contains(err.Error(), "parse config")
}
