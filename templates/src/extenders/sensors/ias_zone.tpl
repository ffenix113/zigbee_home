{{ define "top_level" }}
// For some reason ZB_SCHEDULE_CALLBACK is not defined with current setup,
// and I can't find a necessary include/config to enable it.
// So for now - re-define the callback as another callback.
#ifndef ZBHOME_IAS_ZONE_TOP_LEVEL
#define ZBHOME_IAS_ZONE_TOP_LEVEL

#define ZB_SCHEDULE_CALLBACK ZB_SCHEDULE_APP_CALLBACK
void update_zone_status(zb_bufid_t bufid, zb_uint16_t cb_data) {
	// TODO: Probably this function needs to free bufid somewhere,
	// but I am not sure if it is actually the case.
	// Would need to double-check.

	// Decode values from the argument.
	bool status = cb_data & 1;
	zb_uint8_t endpoint = cb_data >> 1;

	switch (status) {
	case true:
		ZB_ZCL_IAS_ZONE_SET_BITS(bufid, endpoint, ZB_ZCL_IAS_ZONE_ZONE_STATUS_ALARM1);
		break;
	case false:
		ZB_ZCL_IAS_ZONE_CLEAR_BITS(bufid, endpoint, ZB_ZCL_IAS_ZONE_ZONE_STATUS_ALARM1);
		break;
	}
}
#endif
{{ end }}

{{ define "button_changed"}}
	bool button_pressed = button_state & {{ toButtonBitName .Sensor.Button}};
	bool button_changed = has_changed & {{ toButtonBitName .Sensor.Button}};
	if (button_changed) {
		// Pack data.
		// Endpoint can be >127, so we can't pack it and state into single uint8.
		zb_uint16_t data = {{.Endpoint}} << 1 | button_pressed;

		/* Allocate output buffer and send on/off command. */
		zb_ret_t zb_err_code = zb_buf_get_out_delayed_ext(
			update_zone_status, data, 0);
		ZB_ERROR_CHECK(zb_err_code);
	}
{{end}}

{{ define "loop" }} {{end}}
{{ define "main"}} {{ end}}