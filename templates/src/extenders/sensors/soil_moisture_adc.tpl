{{/* The templates are non-empty to force their usage. */}}
{{ define "top_level" }} {{end}}
{{ define "button_changed"}} {{/* button_status = has_button_changed(&{{.Sensor.Pin.Label}}, button_state, has_changed); */}}{{end}}
{{ define "loop"}}
    int32_t sensor_mv;
    int err = zigbee_home_read_adc_mv(&adc_channel_{{.Sensor.ADCPin.Name}}, &sensor_mv);
    if (err) {
        LOG_ERR("Failed to read ADC value from ADC channel {{.Sensor.ADCPin.Name}}");
    }

    // This values are from sensor configuration.
    // Max value is actually minimal moisture,
    // and min value is maximal moisture.
    const zb_uint16_t min_mv_val = {{.Sensor.MaxMoistureMv}};
    const zb_uint16_t max_mv_val = {{.Sensor.MinMoistureMv}};

    // Assume that all ADC moisture sensors will 
    // return smaller values with higher soil moisture content.
    zb_uint16_t percentage;
    if (sensor_mv >= max_mv_val) {
        percentage = 0;
    } else if (sensor_mv <= min_mv_val) {
        percentage = 100;
    } else {
        percentage = 100 - (((sensor_mv-min_mv_val)*100) / (max_mv_val - min_mv_val));
    }

    // Have only 5% increments to not be very noisy,
    // this could change in the future though.
    percentage = ((percentage + 5) / 10) * 10;

    LOG_DBG("ADC soil moisture channel {{.Sensor.ADCPin.Name}}: min_mv: %d, max_mv: %d, current mv: %d, percent: %d", min_mv_val, max_mv_val, sensor_mv, percentage);

    percentage *= 100;
    {{ $cluster := (index .Sensor.Clusters 0) }}
    int status = zb_zcl_set_attr_val({{.Endpoint}},
                                {{ $cluster.ID }},
                                ZB_ZCL_CLUSTER_SERVER_ROLE,
                                ZB_ZCL_ATTR_REL_HUMIDITY_MEASUREMENT_VALUE_ID,
                                (zb_uint8_t*)&percentage,
                                ZB_FALSE);
    if (status) {
        LOG_ERR("Failed to set ZCL attribute for soil moisture sensor: %d", status);
    }
{{end}}
{{ define "main"}} {{end}}