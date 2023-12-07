package config

import (
	"fmt"
	"io"
)

type ConfigValue struct {
	Name          string
	DefaultValue  string
	RequiredValue string

	Dependencies []ConfigValue
}

type Config struct {
	// a simple map, as we store already resolved values.
	values map[string]string
}

func NewValue(name string) ConfigValue {
	return ConfigValue{Name: name}
}

func (v ConfigValue) Default(val string) ConfigValue {
	v.DefaultValue = val
	return v
}

func (v ConfigValue) Required(val string) ConfigValue {
	v.RequiredValue = val
	return v
}

func (v ConfigValue) Value() string {
	if v.RequiredValue != "" {
		return v.RequiredValue
	}

	return v.DefaultValue
}

// TODO: add ConfigValue.Depends(...ConfigValue)

func NewDefaultConfig() *Config {
	return (&Config{}).AddValue(
		CONFIG_LOG,
		CONFIG_SERIAL,
		CONFIG_CONSOLE,
		CONFIG_UART_CONSOLE,
		CONFIG_UART_LINE_CTRL,
		CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT,
		CONFIG_DK_LIBRARY,
		CONFIG_ZIGBEE,
		CONFIG_ZIGBEE_APP_UTILS,
		CONFIG_ZIGBEE_CHANNEL,
		CONFIG_ZIGBEE_ROLE_END_DEVICE,
		CONFIG_CRYPTO,
		CONFIG_CRYPTO_NRF_ECB,
		CONFIG_CRYPTO_INIT_PRIORITY,
		CONFIG_RAM_POWER_DOWN_LIBRARY,
		CONFIG_NET_IPV6,
		CONFIG_NET_IP_ADDR_CHECK,
		CONFIG_NET_UDP,
		CONFIG_ZBOSS_HALT_ON_ASSERT,
		CONFIG_RESET_ON_FATAL_ERROR,
		CONFIG_LOG_BACKEND_UART,
		CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE,
		CONFIG_HEAP_MEM_POOL_SIZE,
	)
}

func (c *Config) AddValue(configValues ...ConfigValue) *Config {
	for _, configValue := range configValues {
		// Only single-level dependencies for now.
		for _, dep := range configValue.Dependencies {
			val, ok := c.values[dep.Name]
			if !ok {
				c.values[dep.Name] = dep.Value()

				continue
			}

			if dep.RequiredValue != "" && val != dep.RequiredValue {
				panic(fmt.Sprintf("config value %q already has value %q, but %q requires it to be %q", dep.Name, val, configValue.Name, dep.RequiredValue))
			}
		}

	}

	return c
}

func (c *Config) Write(w io.StringWriter) error {
	for name, value := range c.values {
		if _, err := w.WriteString(name + "=" + value + "\n"); err != nil {
			return fmt.Errorf("write to writer: %w", err)
		}
	}

	return nil
}
