package templates

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

func typeFromSensor(sensor SensorCtx) string {
	return sensor.Sensor.CPPComponentType()
}

func toButtonBit(btnID string) string {
	return "BUTTON_BIT(" + btnID + ")"
}

func sensorCtx(endpoint int, device *config.Device, sensor sensor.Sensor, extender generator.Extender) SensorCtx {
	return SensorCtx{
		Endpoint: endpoint,
		Device:   device,
		Sensor:   sensor,
		Extender: extender,
	}
}

func clusterCtx(endpoint int, cluster cluster.Cluster) ClusterCtx {
	return ClusterCtx{
		Endpoint: endpoint,
		Cluster:  cluster,
	}
}

func isLast(i, arrLen int) bool {
	return i+1 == arrLen
}

func sum(a, b int) int {
	return a + b
}

func formatHex(val any) (string, error) {
	switch i := val.(type) {
	case uint8, uint16, uint32, uint64, uint,
		int8, int16, int32, int64, int:
		return fmt.Sprintf("%#x", i), nil
	default:
		return "", fmt.Errorf("unknown type to format: %T", val)
	}
}

func formatClusterArgs(args []any) string {
	var result string

	for i, arg := range args {
		if i != 0 {
			result += ","
		}

		if stringer, ok := arg.(fmt.Stringer); ok {
			result += stringer.String()
		} else {
			result += fmt.Sprintf("%v", arg)
		}
	}

	return result
}

// trustCenterKeyToArray will return trust key as C array item values.
func trustCenterKeyToArray(key string) (string, error) {
	key = strings.ReplaceAll(key, ":", "")

	var sb strings.Builder

	for i := 0; i < 16; i++ {
		bytePart := key[i*2 : i*2+2]

		_, err := strconv.ParseUint(bytePart, 16, 16)
		if err != nil {
			return "", fmt.Errorf("parse trust center key part %q: %w", bytePart, err)
		}

		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString("0x" + bytePart)
	}

	return sb.String(), nil
}

func ncsVersionIs(current, another types.Semver) func() bool {
	isSame := current.SameMajorMinor(another)

	return func() bool {
		return isSame
	}
}

// timeToVersion will accept some timestamp
// and will return 4 strings with VERSION-formatted strings.
//
// It is only useful for VERSION file.
//
// TODO: this could return a struct as well.
// Maybe the struct that could format itself to version.
func timeToVersion(now time.Time) [4]string {
	unixSeconds := now.Unix()

	var parts [4]byte
	binary.BigEndian.PutUint32(parts[:], uint32(unixSeconds))

	var formatted [4]string

	for i, part := range parts {
		formatted[i] = strconv.FormatUint(uint64(part), 10)
	}

	return formatted
}
