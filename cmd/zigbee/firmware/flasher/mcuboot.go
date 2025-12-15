package flasher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"gopkg.in/yaml.v3"
	"mynewt.apache.org/newt/util"
	"mynewt.apache.org/newtmgr/newtmgr/cli"
	"mynewt.apache.org/newtmgr/newtmgr/nmutil"
	"mynewt.apache.org/newtmgr/nmxact/nmp"
	"mynewt.apache.org/newtmgr/nmxact/nmxutil"
	"mynewt.apache.org/newtmgr/nmxact/xact"
)

// MCUBoot will flash the board that has mcuboot.
// https://docs.zephyrproject.org/latest/boards/arm/nrf52840dongle_nrf52840/doc/index.html#option-2-using-mcuboot-in-serial-recovery-mode
type MCUBoot struct{}

func (MCUBoot) Flash(ctx context.Context, device *config.Device, workDir string) error {
	// toolchainsPath := device.General.GetToochainsPath()
	// opts := []runner.CmdOpt{
	// 	runner.WithWorkDir(workDir),
	// 	config.WithToolchainPath(toolchainsPath.ToolchainPath, toolchainsPath.SDKPath),
	// }

	logrusLevel := logrus.InfoLevel
	if err := util.Init(logrusLevel, "", util.VERBOSITY_DEFAULT); err != nil {
		return fmt.Errorf("init newt library: %w", err)
	}
	nmxutil.SetLogLevel(logrusLevel)

	cli.OSSpecificInit()

	if err := setupSerial(device); err != nil {
		return fmt.Errorf("setup serial: %w", err)
	}

	if device.Experimental.BLEOTA != nil {
		err := setupBLE(device)
		if err != nil {
			return fmt.Errorf("setup ble: %w", err)
		}

		log.Println("0/2 disabling Zigbee for image upload")
		if err := resetDevice(true); err != nil {
			return fmt.Errorf("reseting before upload: %w", err)
		}
	}

	log.Println("1/2 flashing device")
	if err := flashDevice(workDir); err != nil {
		return fmt.Errorf("flashing device: %w", err)
	}

	log.Println("2/2 reseting device")
	if err := resetDevice(false); err != nil {
		return fmt.Errorf("reseting after upload: %w", err)
	}

	log.Println("Success. Device may take ~2-3 minutes to boot up, please be patient.")

	return nil
}

func setupBLE(cfg *config.Device) error {
	devName, err := extenders.BLEDeviceName(
		cfg.Experimental.BLEOTA.DeviceNamePrefix,
		cfg.General.DeviceName,
		cfg.Experimental.BLEOTA.DeviceName)
	if err != nil {
		return fmt.Errorf("get ble device name: %w", err)
	}

	nmutil.Timeout = 10
	nmutil.MtuOverride = 256

	nmutil.ConnType = "ble"
	nmutil.ConnString = "peer_name=" + devName

	return nil
}

func setupSerial(_ *config.Device) error {
	nmutil.Timeout = 10

	nmutil.ConnType = "serial"
	// FIXME: This would be wrong probably 90% of the time,
	// as device would probaly never be the same,
	// especially on macOs.
	//
	// Would be nice to have ability to iterate over usb devices
	// and find some that match requirements instead.
	nmutil.ConnString = "dev=/dev/ttyACM0,baud=115200"

	return nil
}

func flashDevice(projectDir string) error {
	buildDir, err := getBuildDir(projectDir)
	if err != nil {
		return fmt.Errorf("get build dir: %w", err)
	}

	imageData, err := os.ReadFile(filepath.Join(buildDir, "zephyr", "zephyr.signed.bin"))
	if err != nil {
		return fmt.Errorf("read image file: %w", err)
	}

	session, err := cli.GetSesn()
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}

	c := xact.NewImageUpgradeCmd()
	c.SetTxOptions(nmutil.TxOptions())
	c.Data = imageData
	c.MaxWinSz = xact.IMAGE_UPLOAD_DEF_MAX_WS

	var lastPercent int
	c.ProgressCb = func(c *xact.ImageUploadCmd, r *nmp.ImageUploadRsp) {
		percent := 100 * int(r.Off) / len(imageData)

		if percent == 0 || percent > lastPercent {
			log.Printf("upload progress: %d%% \n", percent)
		}

		lastPercent = percent
	}

	res, err := c.Run(session)
	if err != nil {
		return fmt.Errorf("run command: %w", util.ChildNewtError(err))
	}

	if res.Status() != 0 {
		return fmt.Errorf("return code: %d", res.Status())
	}

	return nil
}

func resetDevice(force bool) error {
	session, err := cli.GetSesn()
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}

	c := xact.NewResetCmd()
	c.SetTxOptions(nmutil.TxOptions())

	c.Force = force

	if _, err := c.Run(session); err != nil {
		return fmt.Errorf("run command: %w", util.ChildNewtError(err))
	}

	return nil
}

func getBuildDir(projectDir string) (string, error) {
	type domains struct {
		Default string `yaml:"default"`
		Domains []struct {
			Name     string `yaml:"name"`
			BuildDir string `yaml:"build_dir"`
		} `yaml:"domains"`
	}

	domainsData, err := os.ReadFile(filepath.Join(projectDir, "build", "domains.yaml"))
	if err != nil {
		return "", fmt.Errorf("read domains.yaml: %w", err)
	}

	var domainsMap domains
	if err := yaml.Unmarshal(domainsData, &domainsMap); err != nil {
		return "", fmt.Errorf("unmarshal domains.yaml: %w", err)
	}

	var buildDir string
	for _, domain := range domainsMap.Domains {
		if domain.Name != domainsMap.Default {
			continue
		}

		buildDir = domain.BuildDir
		break
	}

	if buildDir == "" {
		return "", errors.New("could not find zephyr build directory")
	}

	return buildDir, nil
}
