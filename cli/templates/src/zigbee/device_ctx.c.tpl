{{define "device_ctx"}}

{{ template "additional_types" }}

typedef struct {
	zb_zcl_basic_attrs_ext_t basic_attr;
	// zb_zcl_identify_attrs_t identify_attr;

    {{- range .Clusters }}
    {{ .CAttrType }} {{.CVarName}}_attrs;
    {{- end }}
} zb_device_ctx;

{{ end }}