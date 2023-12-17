{{ template "additional_types" }}

struct zb_device_ctx {
	zb_zcl_basic_attrs_t basic_attr;
	// zb_zcl_identify_attrs_t identify_attr;

    {{ range .Clusters }}
    {{ .CTypeName }} {{.CVarName}}_attrs;
    {{ end }}
};