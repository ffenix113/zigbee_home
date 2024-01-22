package appconfig

import (
	"fmt"
	"io"
	"sort"
)

type Provider interface {
	AppConfig() []ConfigValue
}

type ConfigValue struct {
	Name          string
	DefaultValue  string
	RequiredValue string

	Dependencies []ConfigValue
}

type AppConfig struct {
	values map[string]ConfigValue
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

func (v ConfigValue) Depends(cfgs ...ConfigValue) ConfigValue {
	v = v.Copy()

	for i := range cfgs {
		// Make value required, as we *require* it as a dependency
		cfgs[i].RequiredValue = cfgs[i].Value()
	}

	v.Dependencies = append(v.Dependencies, cfgs...)
	return v
}

func (v ConfigValue) Copy() ConfigValue {
	return ConfigValue{
		Name:          v.Name,
		DefaultValue:  v.DefaultValue,
		RequiredValue: v.RequiredValue,

		Dependencies: append([]ConfigValue(nil), v.Dependencies...),
	}
}

func NewEmptyAppConfig() *AppConfig {
	return &AppConfig{
		values: make(map[string]ConfigValue),
	}
}

func NewDefaultAppConfig(isRouter bool) *AppConfig {
	appConfig := NewEmptyAppConfig().AddValue(
		CONFIG_LOG,
		CONFIG_SERIAL,
		CONFIG_CONSOLE,
		CONFIG_UART_CONSOLE,
		CONFIG_UART_LINE_CTRL,
		CONFIG_USB_DEVICE_STACK,
		CONFIG_USB_DEVICE_INITIALIZE_AT_BOOT,
		CONFIG_DK_LIBRARY,
		CONFIG_ZIGBEE,
		CONFIG_ZIGBEE_APP_UTILS,
		CONFIG_ZIGBEE_CHANNEL,
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

	deviceRole := CONFIG_ZIGBEE_ROLE_END_DEVICE
	if isRouter {
		deviceRole = CONFIG_ZIGBEE_ROLE_ROUTER
	}

	return appConfig.AddValue(deviceRole)
}

func (c *AppConfig) AddValue(configValues ...ConfigValue) *AppConfig {
	for _, configValue := range configValues {
		// Only single-level dependencies for now.
		for _, dep := range configValue.Dependencies {
			val, ok := c.values[dep.Name]
			if !ok {
				c.values[dep.Name] = dep

				continue
			}

			if dep.RequiredValue != "" {
				if val.RequiredValue != "" && val.RequiredValue != dep.RequiredValue {
					panic(fmt.Sprintf("config value %q already has required value %q, but %q requires it to be %q", dep.Name, val.RequiredValue, configValue.Name, dep.RequiredValue))
				}

				c.values[dep.Name] = dep
				continue
			}
		}

		c.values[configValue.Name] = configValue
	}

	return c
}

func (c *AppConfig) WriteTo(w io.StringWriter) error {
	configNames := make([]string, 0, len(c.values))
	for name := range c.values {
		configNames = append(configNames, name)
	}
	sort.Strings(configNames)

	for _, name := range configNames {
		if _, err := w.WriteString(name + "=" + c.values[name].Value() + "\n"); err != nil {
			return fmt.Errorf("write to writer: %w", err)
		}
	}

	return nil
}
