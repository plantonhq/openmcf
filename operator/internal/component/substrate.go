package component

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// substrate names the kind of cluster the operator finds itself on, at the
// granularity that changes our guidance -- not full cloud inventory. The one
// distinction that matters most is EKS vs "AWS but not EKS": the first gets
// IRSA (annotate a ServiceAccount, no keys anywhere), the second gets instance
// profiles or a credentials Secret. Detection is best-effort and read-only;
// an unknown substrate just yields generic guidance, never an error.
type substrate string

const (
	substrateEKS     substrate = "eks"
	substrateAWS     substrate = "aws" // AWS nodes without EKS markers (kOps, RKE2-on-EC2, ...)
	substrateGKE     substrate = "gke"
	substrateAKS     substrate = "aks"
	substrateRKE2    substrate = "rke2" // Rancher RKE2 outside a cloud (datacenter/bare metal)
	substrateK3s     substrate = "k3s"
	substrateUnknown substrate = "unknown"
)

// detectSubstrate classifies the cluster from its nodes. Signals, strongest
// first:
//
//   - spec.providerID scheme ("aws://", "gce://", "azure://") tells the CLOUD
//     -- set by the cloud-controller-manager, present on any integrated node.
//   - kubeletVersion suffix ("-eks-...") and the managed-service node labels
//     tell whether the cloud's MANAGED Kubernetes is running on it.
//
// A single node is enough: mixed-substrate clusters do not exist in practice,
// and the first classifiable node decides.
func detectSubstrate(nodes []corev1.Node) substrate {
	for i := range nodes {
		if s := classifyNode(&nodes[i]); s != substrateUnknown {
			return s
		}
	}
	return substrateUnknown
}

func classifyNode(node *corev1.Node) substrate {
	providerID := node.Spec.ProviderID
	kubelet := node.Status.NodeInfo.KubeletVersion
	labels := node.Labels

	switch {
	case strings.HasPrefix(providerID, "aws://"):
		// EKS marks nodes two ways; either suffices. Distros that merely RUN
		// on EC2 (kOps, RKE2, kubeadm) carry the aws:// providerID but
		// neither marker.
		if strings.Contains(kubelet, "-eks-") || labels["eks.amazonaws.com/nodegroup"] != "" {
			return substrateEKS
		}
		return substrateAWS
	case strings.HasPrefix(providerID, "gce://"):
		if strings.Contains(kubelet, "-gke.") || labels["cloud.google.com/gke-nodepool"] != "" {
			return substrateGKE
		}
		return substrateUnknown
	case strings.HasPrefix(providerID, "azure://"):
		if labels["kubernetes.azure.com/cluster"] != "" || labels["kubernetes.azure.com/agentpool"] != "" {
			return substrateAKS
		}
		return substrateUnknown
	default:
		// No cloud providerID: a datacenter/bare-metal cluster. Rancher's
		// distributions stamp the kubelet version ("v1.31.4+rke2r1",
		// "v1.31.4+k3s1") -- naming them beats "unknown" because their
		// conventions differ from clouds (node-served ingress, no cloud LB).
		// Cloud-hosted RKE2 stays classified by its CLOUD above: credential
		// and storage guidance follow the cloud, not the distribution.
		// Mirrored by the guided installer's detector (selfhosted-install).
		switch {
		case strings.Contains(kubelet, "+rke2"):
			return substrateRKE2
		case strings.Contains(kubelet, "+k3s"):
			return substrateK3s
		}
		return substrateUnknown
	}
}
