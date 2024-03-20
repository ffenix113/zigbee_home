package appconfig

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

type Provider interface {
	AppConfig() []ConfigValue
}

type ConfigValue struct {
	Name          string
	DefaultValue  string
	RequiredValue string
	QuotedValue   bool

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

func (v ConfigValue) Quoted() ConfigValue {
	v.QuotedValue = true
	return v
}

func (v ConfigValue) Value() string {
	returnValue := v.DefaultValue

	if v.RequiredValue != "" {
		returnValue = v.RequiredValue
	}

	if v.QuotedValue {
		returnValue = `"` + returnValue + `"`
	}

	return returnValue
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

type DefaultAppConfigOptions struct {
	IsRouter       bool
	ZigbeeChannels []int
}

func NewDefaultAppConfig(opts DefaultAppConfigOptions) (*AppConfig, error) {
	appConfig := NewEmptyAppConfig().AddValue(
		CONFIG_DK_LIBRARY,
		CONFIG_ZIGBEE,
		CONFIG_ZIGBEE_APP_UTILS,
		CONFIG_ZIGBEE_CHANNEL_MASK,
		CONFIG_ZIGBEE_CHANNEL_SELECTION_MODE_MULTI,
		CONFIG_CRYPTO,
		CONFIG_CRYPTO_NRF_ECB,
		CONFIG_CRYPTO_INIT_PRIORITY,
		CONFIG_RAM_POWER_DOWN_LIBRARY,
		CONFIG_NET_IPV6,
		CONFIG_NET_IP_ADDR_CHECK,
		CONFIG_NET_UDP,
		CONFIG_ZBOSS_HALT_ON_ASSERT,
		CONFIG_RESET_ON_FATAL_ERROR,
		CONFIG_SYSTEM_WORKQUEUE_STACK_SIZE,
		CONFIG_HEAP_MEM_POOL_SIZE,
	)

	deviceRole := CONFIG_ZIGBEE_ROLE_END_DEVICE
	if opts.IsRouter {
		deviceRole = CONFIG_ZIGBEE_ROLE_ROUTER
	}

	appConfig = appConfig.AddValue(deviceRole)

	if len(opts.ZigbeeChannels) != 0 {
		channel := int32(0)
		for _, chann := range opts.ZigbeeChannels {
			if chann < 11 || chann > 26 {
				return nil, fmt.Errorf("zigbee channels must be in range [11, 26], but have %d", chann)
			}

			channel |= 1 << chann
		}

		appConfig = appConfig.AddValue(CONFIG_ZIGBEE_CHANNEL_MASK.Required("0x" + strconv.FormatInt(int64(channel), 16)))
	}

	return appConfig, nil
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

		if val, ok := c.values[configValue.Name]; ok {
			if val.RequiredValue != "" &&
				configValue.RequiredValue != "" &&
				configValue.RequiredValue != val.RequiredValue {
				panic(fmt.Sprintf("config value %q already has required value %q, but new added value requires it to be %q", val.Name, val.RequiredValue, configValue.RequiredValue))
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
