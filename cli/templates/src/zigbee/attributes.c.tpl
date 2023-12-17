{{ define "additional_types" }}

// Types for some clusters
typedef struct {
	zb_int16_t measure_value;
	zb_int16_t min_measure_value;
	zb_int16_t max_measure_value;
	zb_uint16_t tolerance;
} zb_zcl_pressure_measurement_attrs_t;

typedef struct {
	zb_int16_t measure_value;
	zb_int16_t min_measure_value;
	zb_int16_t max_measure_value;
} zb_zcl_humidity_measurement_attrs_t;

//
{{end}}