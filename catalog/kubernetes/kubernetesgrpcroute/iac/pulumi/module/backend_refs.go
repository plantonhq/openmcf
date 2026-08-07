package module

import (
	kubernetesgrpcroutev1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesgrpcroute/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildBackendRefs maps the Planton backend references (flattened group / kind /
// name / namespace / port / weight plus optional per-backend filters) onto the
// typed crd2pulumi backendRefs array.
//
// Each backend's name is a KubernetesService foreign key resolved to its
// literal value before the module runs, so GetValue() returns the final name.
func buildBackendRefs(backendRefs []*kubernetesgrpcroutev1alpha1.KubernetesGrpcRouteBackendRef) gatewayv1.GRPCRouteSpecRulesBackendRefsArray {
	arr := gatewayv1.GRPCRouteSpecRulesBackendRefsArray{}
	for _, b := range backendRefs {
		args := gatewayv1.GRPCRouteSpecRulesBackendRefsArgs{
			Name: pulumi.String(b.GetName().GetValue()),
		}
		if group := b.GetGroup(); group != "" {
			args.Group = pulumi.String(group)
		}
		if kind := b.GetKind(); kind != "" {
			args.Kind = pulumi.String(kind)
		}
		if namespace := b.GetNamespace(); namespace != "" {
			args.Namespace = pulumi.String(namespace)
		}
		if b.Port != nil {
			args.Port = pulumi.Int(int(b.GetPort()))
		}
		if b.Weight != nil {
			args.Weight = pulumi.Int(int(b.GetWeight()))
		}
		if filters := b.GetFilters(); len(filters) > 0 {
			args.Filters = buildBackendFilters(filters)
		}
		arr = append(arr, args)
	}
	return arr
}

// buildBackendFilters maps the same Planton KubernetesGrpcRouteFilter message
// onto the backend-ref-level crd2pulumi filter type tree. This mirrors
// buildRuleFilters in filters.go; the two trees are kept separate because
// crd2pulumi generates distinct (GRPCRouteSpecRulesFilters* vs
// GRPCRouteSpecRulesBackendRefsFilters*) Go types for the same JSON shape.
func buildBackendFilters(filters []*kubernetesgrpcroutev1alpha1.KubernetesGrpcRouteFilter) gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersArray {
	arr := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersArray{}
	for _, f := range filters {
		args := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersArgs{
			Type: pulumi.String(f.GetType()),
		}
		if v := f.GetRequestHeaderModifier(); v != nil {
			args.RequestHeaderModifier = buildBackendRequestHeaderModifier(v)
		}
		if v := f.GetResponseHeaderModifier(); v != nil {
			args.ResponseHeaderModifier = buildBackendResponseHeaderModifier(v)
		}
		if v := f.GetRequestMirror(); v != nil {
			args.RequestMirror = buildBackendRequestMirror(v)
		}
		if v := f.GetExtensionRef(); v != nil {
			args.ExtensionRef = gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersExtensionRefArgs{
				Group: pulumi.String(v.GetGroup()),
				Kind:  pulumi.String(v.GetKind()),
				Name:  pulumi.String(v.GetName()),
			}
		}
		arr = append(arr, args)
	}
	return arr
}

func buildBackendRequestHeaderModifier(m *kubernetesgrpcroutev1alpha1.KubernetesGrpcRouteHeaderFilter) gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestHeaderModifierArgs {
	args := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestHeaderModifierArgs{}
	if set := m.GetSet(); len(set) > 0 {
		setArr := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestHeaderModifierSetArray{}
		for _, h := range set {
			setArr = append(setArr, gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestHeaderModifierSetArgs{
				Name:  pulumi.String(h.GetName()),
				Value: pulumi.String(h.GetValue()),
			})
		}
		args.Set = setArr
	}
	if add := m.GetAdd(); len(add) > 0 {
		addArr := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestHeaderModifierAddArray{}
		for _, h := range add {
			addArr = append(addArr, gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestHeaderModifierAddArgs{
				Name:  pulumi.String(h.GetName()),
				Value: pulumi.String(h.GetValue()),
			})
		}
		args.Add = addArr
	}
	if remove := m.GetRemove(); len(remove) > 0 {
		args.Remove = pulumi.ToStringArray(remove)
	}
	return args
}

func buildBackendResponseHeaderModifier(m *kubernetesgrpcroutev1alpha1.KubernetesGrpcRouteHeaderFilter) gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersResponseHeaderModifierArgs {
	args := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersResponseHeaderModifierArgs{}
	if set := m.GetSet(); len(set) > 0 {
		setArr := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersResponseHeaderModifierSetArray{}
		for _, h := range set {
			setArr = append(setArr, gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersResponseHeaderModifierSetArgs{
				Name:  pulumi.String(h.GetName()),
				Value: pulumi.String(h.GetValue()),
			})
		}
		args.Set = setArr
	}
	if add := m.GetAdd(); len(add) > 0 {
		addArr := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersResponseHeaderModifierAddArray{}
		for _, h := range add {
			addArr = append(addArr, gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersResponseHeaderModifierAddArgs{
				Name:  pulumi.String(h.GetName()),
				Value: pulumi.String(h.GetValue()),
			})
		}
		args.Add = addArr
	}
	if remove := m.GetRemove(); len(remove) > 0 {
		args.Remove = pulumi.ToStringArray(remove)
	}
	return args
}

func buildBackendRequestMirror(m *kubernetesgrpcroutev1alpha1.KubernetesGrpcRouteRequestMirrorFilter) gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestMirrorArgs {
	args := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestMirrorArgs{}
	if ref := m.GetBackendRef(); ref != nil {
		backendRef := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestMirrorBackendRefArgs{
			Name: pulumi.String(ref.GetName().GetValue()),
		}
		if group := ref.GetGroup(); group != "" {
			backendRef.Group = pulumi.String(group)
		}
		if kind := ref.GetKind(); kind != "" {
			backendRef.Kind = pulumi.String(kind)
		}
		if namespace := ref.GetNamespace(); namespace != "" {
			backendRef.Namespace = pulumi.String(namespace)
		}
		if ref.Port != nil {
			backendRef.Port = pulumi.Int(int(ref.GetPort()))
		}
		args.BackendRef = backendRef
	}
	if m.Percent != nil {
		args.Percent = pulumi.Int(int(m.GetPercent()))
	}
	if fraction := m.GetFraction(); fraction != nil {
		fractionArgs := gatewayv1.GRPCRouteSpecRulesBackendRefsFiltersRequestMirrorFractionArgs{
			Numerator: pulumi.Int(int(fraction.GetNumerator())),
		}
		if fraction.Denominator != nil {
			fractionArgs.Denominator = pulumi.Int(int(fraction.GetDenominator()))
		}
		args.Fraction = fractionArgs
	}
	return args
}
