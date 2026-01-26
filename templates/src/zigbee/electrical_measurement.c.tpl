{{ define "electrical_measurement_defines"}}
#include "cluster_electrical_measurement.hpp"
{{ end }}

{{ define "electrical_measurement_attr_list" }}
ZB_ZCL_DECLARE_ELECTRICAL_MEASUREMENT_ATTRIB_LIST(
    {{.Cluster.CVarName}}_{{.Endpoint}}_attr_list, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.measurement_type, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.dc_voltage, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.dc_current, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.dc_power, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.dc_voltage_divisor, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.dc_current_divisor, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.dc_power_divisor, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.one, 
    &dev_ctx.{{.Cluster.CVarName}}_{{.Endpoint}}_attrs.uone);

{{ end }}