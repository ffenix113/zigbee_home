{{/* The templates are non-empty to force their usage. */}}
{{ define "top_level" }} {{end}}
{{ define "button_changed"}} {{/* button_status = has_button_changed(&{{.Sensor.Pin.Label}}, button_state, has_changed); */}}{{end}}
{{ define "loop"}}

    int32_t batt_mv;
    int err = zigbee_home_read_adc_mv(&adc_channel_{{.Sensor.ADCPin.Name}}, &batt_mv);
    if (err) {
        LOG_ERR("Failed to read ADC value from ADC channel {{.Sensor.ADCPin.Name}}");
    }
    {{ $cluster := (index .Sensor.Clusters 0) }}
    zb_uint8_t batt_mv_divided = batt_mv / 100;
    zb_zcl_status_t status = zb_zcl_set_attr_val({{.Endpoint}},
                                {{ $cluster.ID }},
                                ZB_ZCL_CLUSTER_SERVER_ROLE,
                                ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_VOLTAGE_ID,
                                &(batt_mv_divided),
                                ZB_FALSE);
    if (status) {
        LOG_ERR("Failed to set ZCL attribute for battery voltage sensor: %d", status);
    }
    {/* Endpoint subtracts 1 because seems that there is a bug 
        in how endpoints are calculated for attribute structs.
     */}
    zb_uint16_t rated_voltage = dev_ctx.{{$cluster.CVarName}}_{{sum .Endpoint -1}}_attrs.rated_voltage * 100;
    zb_uint16_t min_threshold = dev_ctx.{{$cluster.CVarName}}_{{sum .Endpoint -1}}_attrs.voltage_min_threshold * 100;

    zb_uint8_t percentage;    
    if (batt_mv >= rated_voltage) {
        percentage = 200;
    } else if (batt_mv <= min_threshold) {
        percentage = 0;
    } else {
        percentage = (((batt_mv-min_threshold)*200) / (rated_voltage - min_threshold));
    }

    // Have only 10% increments
    percentage = ((percentage + 10) / 20) * 20;

    LOG_DBG("ADC battery channel {{.Sensor.ADCPin.Name}}: rated: %d, threshold: %d, current mv: %d, percent(x2): %d", rated_voltage, min_threshold, batt_mv, percentage);
    status = zb_zcl_set_attr_val({{.Endpoint}},
                                {{ $cluster.ID }},
                                ZB_ZCL_CLUSTER_SERVER_ROLE,
                                ZB_ZCL_ATTR_POWER_CONFIG_BATTERY_PERCENTAGE_REMAINING_ID,
                                &percentage,
                                ZB_FALSE);
    if (status) {
        LOG_ERR("Failed to set ZCL attribute for battery percentage sensor: %d", status);
    }
{{end}}
{{ define "main"}} {{end}}