package extenders

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/generator"
)

const bleNamePrefix = "ZBHome"

// BLEOTA provides a BLE OTA for device, if device configuration allows it.
//
// Zigbee OTA is not supported (yet?), as it requires additional moving parts
// that this project cannot configure/provide. Same for Matter OTA.
//
// Note: This option will increase power usage for device,
// as it will enable Bluetooth LE and will constantly advertise itself.
// It is not possible to disable BLE advertising and enable it on-request,
// for now, though it may arrive later with next iterations.
// A suggestion would be to have constant power source for a board, like a charger.
// This is experimental feature, it may not be perfect.
//
// Device must have a separate partition, which should be half of the available size,
// or, at least, around 450+ kb per partition. If this requirement is not met -
// build will fail. Currently, it is not possible to fix it via configuration,
// only with manual overlay file.
// In the future zigbee_home will force partitions to support this for all devices.
//
// It also requires the device to have MCUBoot as a first stage bootloader.
// This means that device must be flashed with SWD when
// flashing this firmware for the first time after configuring OTA,
// and not bootloader-specific way.
// Otherwise MCUBoot will not have keys that are used to verify the images.
//
// Generally it is easily possible to also move back to original bootloader
// provided by manufacturer by flashing it via SWD.
type BLEOTA struct {
	// DeviceNamePrefix allows to change default name prefix.
	//
	// This is optional, and can be set to empty string to remove the prefix.
	DeviceNamePrefix *string `yaml:"device_name_prefix"`
	// DeviceName is optional device BLE name suffix.
	//
	// This is so updater will know which device to connect to and update.
	//
	// If this option is not provided - the device name will be used.
	// But be warned that this may easily result in weird behavior when
	// doing OTA, as multiple devices could have the same name.
	//
	// It should be unique per device, so if there is 2 devices named "Sensor" -
	// they may not correctly be updated and some weird behavior may be observed
	// when trying to update the device.
	//
	// If `Key` is set - it will prevent wrong device being updated
	// with incorrect firmware, but user still may experience
	// weird behavior like device randomly fails to update.
	DeviceName string `yaml:"device_name"`

	// Key is optional value to generate private key for this device.
	//
	// This key must be unique per device.
	//
	// The length of the key is not limited, but should not
	// be less than ~5-6 characterst. Secure
	//
	// While it is optional, it is strongly recommended to set it.
	// Key is needed so another board will not be updated accidentaly
	// and to prevent some "rogue" update would not compromise the device.
	// So it just acts as second layer of device & firmware verification.
	// It may also be used in the future for authenticity check for some operations
	// (including OTA).
	//
	// Note: Key cannot be changed after firmware with this key was already flashed,
	// unless it is flashed via SWD/JLink or other method that overwrites bootloader as well.
	// Otherwise the bootloader will not accept the firmware with changed key.
	//
	// FIXME: Should not be here, should be a part of MCUBoot configuration.
	Key string
}

// NewBLEOTA will create an extender that would allow board to be updated
// via BLE OTA.
func NewBLEOTA(workdir string, deviceName string, ota BLEOTA) (generator.Extender, error) {
	bleDeviceName, err := BLEDeviceName(ota.DeviceNamePrefix, deviceName, ota.DeviceName)
	if err != nil {
		return nil, fmt.Errorf("get ble device name: %w", err)
	}

	// FIXME: Should not be here, should be a part of MCUBoot configuration.
	parsedMCUBootKey, err := getMCUBootKey(ota.Key)
	if err != nil {
		return nil, fmt.Errorf("ble ota key: %w", err)
	}

	otaExtender := generator.SimpleExtender{
		Name: "BLE_OTA",
		Config: []appconfig.ConfigValue{
			appconfig.CONFIG_BT,
			appconfig.CONFIG_BT_PERIPHERAL,
			appconfig.CONFIG_BT_DEVICE_NAME.Default(bleDeviceName),
			appconfig.CONFIG_NCS_SAMPLE_MCUMGR_BT_OTA_DFU,
			// Configuration to be able to manage firmware
			// and handle required DFU events
			appconfig.CONFIG_BOOTLOADER_MCUBOOT,
			appconfig.CONFIG_MCUMGR_GRP_IMG_STATUS_HOOKS,
			appconfig.CONFIG_MCUMGR_GRP_OS_RESET_HOOK,
			appconfig.CONFIG_MCUMGR_MGMT_NOTIFICATION_HOOKS,
		},
		Sysbuild: []appconfig.ConfigValue{
			appconfig.SB_CONFIG_BOOTLOADER_MCUBOOT.Required(appconfig.Yes),
		},
		FilesToWrite: []generator.WriteFile{
			{
				TemplateName: "ble_ota.hpp",
				FileName:     "ble_ota.hpp",
			},
			{
				TemplateName: "ble_ota.cpp",
				FileName:     "ble_ota.cpp",
			},
		},
		IncludeHeaders: []string{"ble_ota.hpp"},
	}

	// FIXME: This should be part of MCUBoot setup.
	if parsedMCUBootKey != nil {
		const mcubootKeyFile = "mcuboot_key.pem"

		keyPem, err := marshalPrivateKeyToPEM(parsedMCUBootKey)
		if err != nil {
			return nil, fmt.Errorf("encode ota key: %w", err)
		}

		otaExtender.Config = append(otaExtender.Config,
			appconfig.CONFIG_MCUBOOT_SIGNATURE_KEY_FILE,
		)

		otaExtender.Sysbuild = append(otaExtender.Sysbuild,
			appconfig.SB_CONFIG_BOOT_SIGNATURE_KEY_FILE,
		)

		otaExtender.FilesToWrite = append(otaExtender.FilesToWrite,
			generator.WriteFile{
				TemplateName: mcubootKeyFile,
				FileName:     filepath.Join("..", mcubootKeyFile),
				AdditionalContext: map[string]any{
					"keyPem": keyPem,
				},
			},
		)
	}

	return otaExtender, nil
}

func BLEDeviceName(blePrefix *string, deviceName, otaDeviceName string) (string, error) {
	selectedName := strings.TrimSpace(otaDeviceName)
	if selectedName == "" {
		selectedName = strings.TrimSpace(deviceName)
	}

	if selectedName == "" {
		return "", errors.New("either device name or OTA device name must be provided")
	}
	// Maybe user does not want to have prefix(i.e. for privacy)
	prefix := bleNamePrefix
	if blePrefix != nil {
		prefix = strings.TrimSpace(*blePrefix)
	}

	if prefix != "" {
		prefix += " "
	}

	return prefix + selectedName, nil
}

// FIXME: This should be moved to some MCUBoot setup, not here.
// Actually, even `Key` from BLEOTA should be in the MCUBoot setup.
func getMCUBootKey(keySeed string) (ed25519.PrivateKey, error) {
	if keySeed == "" {
		return nil, nil
	}
	// Reasonable and reproducable way to get the 32 byte key from the input.
	signKey := sha256.Sum256([]byte(keySeed))

	// We can do this knowing that this will not panic,
	// as sha256 is exactly 32 bytes, which is what this function wants.
	return ed25519.NewKeyFromSeed(signKey[:]), nil
}

// FIXME: This should be moved to some MCUBoot setup, not here.
func marshalPrivateKeyToPEM(pk ed25519.PrivateKey) (string, error) {
	bts, err := x509.MarshalPKCS8PrivateKey(pk)
	if err != nil {
		return "", fmt.Errorf("marshal private key: %w", err)
	}

	block := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bts,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &block); err != nil {
		return "", fmt.Errorf("write private key: %w", err)
	}

	return buf.String(), nil
}
