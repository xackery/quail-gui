package config

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	instance *Config
)

// Config represents a configuration parse
type Config struct {
	baseName     string
	IsVirtualWld bool
}

// New creates a new configuration
func New(ctx context.Context, baseName string) (*Config, error) {
	var f *os.File
	cfg := &Config{
		baseName: baseName,
	}
	instance = cfg
	path := baseName + ".ini"

	isNewConfig := false
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("config info: %w", err)
		}
		f, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("create %s.ini: %w", baseName, err)
		}
		fi, err = os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("new config info: %w", err)
		}
		isNewConfig = true
	}
	if !isNewConfig {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open config: %w", err)
		}
	}

	defer f.Close()
	if fi.IsDir() {
		return nil, fmt.Errorf("%s.ini is a directory, should be a file", baseName)
	}

	if isNewConfig {
		cfg = &Config{
			IsVirtualWld: true,
		}
		err = cfg.Save()
		if err != nil {
			return nil, fmt.Errorf("save config: %w", err)
		}
		return cfg, nil
	}

	err = decode(f, cfg)
	if err != nil {
		return nil, fmt.Errorf("decode %s.ini: %w", baseName, err)
	}

	return cfg, nil
}

func Instance() *Config {
	return instance
}

func ByKey(key string) (string, error) {
	if instance == nil {
		return "", fmt.Errorf("config not loaded")
	}
	switch key {
	case "is_virtual_wld":
		return fmt.Sprintf("%v", instance.IsVirtualWld), nil

	}
	return "", fmt.Errorf("unknown key: %s", key)
}

func SetByKey(key, value string) error {
	if instance == nil {
		return fmt.Errorf("config not loaded")
	}
	switch key {
	case "is_virtual_wld":
		instance.IsVirtualWld = value == "true"
	}
	return fmt.Errorf("unknown key: %s", key)
}

// Verify returns an error if configuration appears off
func (c *Config) Verify() error {

	return nil
}

func decode(r io.Reader, cfg *Config) error {
	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		switch key {
		case "is_virtual_wld":
			if strings.ToLower(value) == "true" {
				cfg.IsVirtualWld = true
			}
		}
	}
	return nil
}

// Save saves the config
func (c *Config) Save() error {

	fi, err := os.Stat(c.baseName + ".ini")
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s.ini:  %w", c.baseName, err)
		}
		w, err := os.Create(c.baseName + ".ini")
		if err != nil {
			return fmt.Errorf("create %s.ini: %w", c.baseName, err)
		}
		w.Close()
	}
	if fi != nil && fi.IsDir() {
		return fmt.Errorf("dirCheck %s.ini: is a directory", c.baseName)
	}

	r, err := os.Open(c.baseName + ".ini")
	if err != nil {
		return err
	}
	defer r.Close()

	tmpConfig := &Config{}

	out := ""
	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "#") {
			out += line + "\n"
			continue
		}
		if !strings.Contains(line, "=") {
			out += line + "\n"
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {

		case "is_virtual_wld":
			if tmpConfig.IsVirtualWld {
				continue
			}
			if c.IsVirtualWld {
				value = "true"
			} else {
				value = "false"
			}
			tmpConfig.IsVirtualWld = true
		}
		line = fmt.Sprintf("%s = %s", key, value)
		out += line + "\n"
	}

	if !tmpConfig.IsVirtualWld {
		if c.IsVirtualWld {
			out += "is_virtual_wld = true\n"
		} else {
			out += "is_virtual_wld = false\n"
		}
	}

	err = os.WriteFile(c.baseName+".ini", []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
