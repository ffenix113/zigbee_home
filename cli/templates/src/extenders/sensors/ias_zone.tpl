{{ define "top_level" }}
// Button will be defined by Buttons extender.

// For some reason ZB_SCHEDULE_CALLBACK is not defined with current setup,
// and I can't find a necessary include/config to enable it.
// So for now - re-define the callback as another callback.
#define ZB_SCHEDULE_CALLBACK ZB_SCHEDULE_APP_CALLBACK
void update_zone_status(zb_bufid_t bufid, bool status) {
	// TODO: Probably this function needs to free bufid somewhere,
	// but I am not sure if it is actually the case.
	// Would need to double-check.
	switch (status) {
	case true:
		ZB_ZCL_IAS_ZONE_SET_BITS(bufid, 2, ZB_ZCL_IAS_ZONE_ZONE_STATUS_ALARM1);
		break;
	case false:
		ZB_ZCL_IAS_ZONE_CLEAR_BITS(bufid, 2, ZB_ZCL_IAS_ZONE_ZONE_STATUS_ALARM1);
		break;
	}
}
{{ end }}

{{ define "button_changed"}}
	button_status = has_button_changed(&{{.Sensor.Pin.Label}}, button_state, has_changed);
	if (button_status.has_changed) {
		/* Allocate output buffer and send on/off command. */
		zb_ret_t zb_err_code = zb_buf_get_out_delayed_ext(
			update_zone_status, button_status.state, 0);
		ZB_ERROR_CHECK(zb_err_code);
	}
{{end}}

{{ define "loop" }} {{end}}
{{ define "main"}} {{ end}}