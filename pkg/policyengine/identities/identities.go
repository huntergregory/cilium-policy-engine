package identities

import (
	"sync"

	"github.com/cilium/cilium/pkg/identity"
	"github.com/cilium/cilium/pkg/k8s"
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/labels"
	"github.com/cilium/cilium/pkg/policy"
)

func CreateIdentities(repo *policy.Repository, pods []*slim_corev1.Pod) map[string]*identity.Identity {
	// create identities from pods
	lblsToIdentity := make(map[string]*identity.Identity)
	identityToLabels := make(map[identity.NumericIdentity]labels.LabelArray)
	// must start high enough so as not to conflict with reserved identities (e.g. 1 = host)
	nextID := uint64(5000)
	for _, pod := range pods {
		lbls := ciliumEndpointLabels(pod)
		if _, ok := lblsToIdentity[lbls.String()]; ok {
			continue
		}
		id := identity.NewIdentity(identity.NumericIdentity(nextID), lbls)
		identityToLabels[id.ID] = id.LabelArray
		lblsToIdentity[lbls.String()] = id
		nextID++
	}

	// add identities to selector cache
	wg := &sync.WaitGroup{}
	repo.GetSelectorCache().UpdateIdentities(identityToLabels, nil, wg)

	return lblsToIdentity
}

func ciliumEndpointLabels(pod *slim_corev1.Pod) labels.Labels {
	ns := &slim_corev1.Namespace{
		ObjectMeta: slim_metav1.ObjectMeta{
			Name: pod.Namespace,
			Labels: map[string]string{
				"kubernetes.io/metadata.name": pod.Namespace,
			},
		},
	}
	_, ciliumLabels := k8s.GetPodMetadata(ns, pod)
	lbls := make(labels.Labels, len(ciliumLabels))
	for k, v := range ciliumLabels {
		lbls[k] = labels.Label{
			Key:    k,
			Value:  v,
			Source: labels.LabelSourceK8s,
		}
	}
	return lbls
}
