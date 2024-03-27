{{ define "additional_types" }}

// Types for some clusters
typedef struct {
	zb_int16_t current_temperature;
} zb_zcl_device_temperature_config_attrs_t;

typedef struct {
	zb_int16_t measure_value;
	zb_int16_t min_measure_value;
	zb_int16_t max_measure_value;
	zb_uint16_t tolerance;
} zb_zcl_pressure_measurement_attrs_t;

typedef struct {
	zb_uint32_t measure_value;
	zb_uint32_t min_measure_value;
	zb_uint32_t max_measure_value;
	zb_uint32_t tolerance;
} zb_zcl_measurement_type_single_attrs_t;

typedef struct {
	zb_int16_t measure_value;
	zb_int16_t min_measure_value;
	zb_int16_t max_measure_value;
} zb_zcl_humidity_measurement_attrs_t;

typedef struct {
	zb_int8_t zone_state;
	zb_int16_t zone_type;
	zb_int16_t zone_status;
	zb_uint8_t zone_id;
	zb_int64_t ias_cie_address;
	zb_int8_t cie_short_addr;
	zb_int16_t cie_ep;
	zb_uint8_t number_of_zone_sens_levels_supported;
	zb_uint8_t current_zone_sens_level;
} zb_zcl_ias_zone_attrs_t;

typedef struct {
	zb_int8_t voltage;
	zb_int16_t size;
	zb_int16_t quantity;
	zb_uint8_t rated_voltage;
	zb_int64_t alarm_mask;
	zb_int8_t voltage_min_threshold;
} zb_zcl_power_config_attrs_t;

//
{{end}}