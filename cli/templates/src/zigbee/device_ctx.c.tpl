struct zb_device_ctx {
	zb_zcl_basic_attrs_t basic_attr;
	zb_zcl_identify_attrs_t identify_attr;

    {{ range .Clusters }}
    {{ clusterToCType . }} {{clusterToVar . }};
    {{ end }}

	zb_zcl_temp_measurement_attrs_t temp_attrs;
	struct zb_zcl_pressure_measurement_attrs_t pres_attrs;
	struct zb_zcl_humidity_measurement_attrs_t humidity_attrs;
	zb_zcl_on_off_attrs_t on_off_attr;
};