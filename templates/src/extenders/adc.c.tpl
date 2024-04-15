#include <stdint.h>
#include <zephyr/drivers/adc.h>
#include <zephyr/logging/log.h>

LOG_MODULE_DECLARE(app, LOG_LEVEL_INF);

int zigbee_home_read_adc_mv(const struct adc_dt_spec *spec, int32_t *valp) {
	int err;
	uint16_t buf;
	struct adc_sequence sequence = {
		.buffer = &buf,
		/* buffer size in bytes, not number of samples */
		.buffer_size = sizeof(buf),
	};

	(void)adc_sequence_init_dt(spec, &sequence);
	err = adc_read(spec->dev, &sequence);
	if (err < 0) {
		LOG_DBG("ADC %s@%d: Could not read (%d)\n", spec->dev->name, spec->channel_id, err);
		return err;
	}

	/*
	* If using differential mode, the 16 bit value
	* in the ADC sample buffer should be a signed 2's
	* complement value.
	*/
	int32_t val_mv;
	if (spec->channel_cfg.differential) {
		val_mv = (int32_t)((int16_t)buf);
	} else {
		val_mv = (int32_t)buf;
	}

	LOG_DBG("ADC %s@%d raw value: %d", spec->dev->name, spec->channel_id, val_mv);
	err = adc_raw_to_millivolts_dt(spec, &val_mv);
	/* conversion to mV may not be supported, skip if not */
	if (err < 0) {
		LOG_DBG("  (value in mV not available)");
		return err;
	}

	LOG_DBG("ADC %s@%d mv value: %d", spec->dev->name, spec->channel_id, val_mv);

	*valp = val_mv;
	return 0;
}