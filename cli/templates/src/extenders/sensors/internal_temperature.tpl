{{ define "main"}}
nrfx_temp_config_t config = NRFX_TEMP_DEFAULT_CONFIG;
int err = nrfx_temp_init(&config, NULL);
if (err != NRFX_SUCCESS) 
{
    LOG_ERR("error initing device temp sensor");
    return err;
}
{{ end}}

{{ define "loop"}}
int res = nrfx_temp_measure();
if (res != NRFX_SUCCESS) {
    LOG_ERR("measure internal temperature failed");
} else {
    int32_t device_temp = nrfx_temp_result_get();
    int32_t real_device_temp = nrfx_temp_calculate(device_temp);

    /* Set ZCL attribute */
    {{ $cluster := (index .Sensor.Clusters 0) }}
    zb_zcl_status_t status = zb_zcl_set_attr_val({{.Endpoint}},
                                {{ $cluster.ID }},
                                ZB_ZCL_CLUSTER_SERVER_ROLE,
                                0x000, // CurrentTemperature attr. Hardcoding for now. Need better approach.
                                (zb_uint8_t *)&real_device_temp,
                                ZB_FALSE);
    if (status) {
        LOG_ERR("Failed to set ZCL attribute for internal temp: %d", status);
    }
}
{{ end}}