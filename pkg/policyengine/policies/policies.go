package policies

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/cilium/cilium/pkg/identity"
	"github.com/cilium/cilium/pkg/k8s"
	slim_networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/networking/v1"
	"github.com/cilium/cilium/pkg/policy"
	"github.com/cilium/cilium/pkg/policy/api"
)

var rev uint64

// NOTE: cidr and namedport not yet supported

func AddPolicies(repo *policy.Repository, netPols []*slim_networkingv1.NetworkPolicy) {
	// parse network policies into api.Rules and add to policy repository
	rules := make(api.Rules, 0)
	for _, netPol := range netPols {
		r, err := k8s.ParseNetworkPolicy(netPol)
		if err != nil {
			logrus.WithError(err).Error("failed to parse NetworkPolicy")
			continue
		}
		rules = append(rules, r...)
	}
	repo.MustAddList(rules)
}

func CountEntries(repo *policy.Repository, lblsToIdentity map[string]*identity.Identity) map[string]int {
	identityEntries := make(map[string]int)
	// count policy map entries for each pod
	for lbls, id := range lblsToIdentity {
		// convert network policies to an endpoint policy for this identity
		endpointPolicy, err := toEndpointPolicy(repo, &fakeEndpoint{}, id)
		if err != nil {
			logrus.WithField("podLabels", lbls).WithError(err).Error("failed to convert Network Policy to Endpoint Policy")
			continue
		}

		// apply policy map changes on the endpoint policy
		closer, _ := endpointPolicy.ConsumeMapChanges()
		closer()

		// get number of policy map entries
		identityEntries[lbls] = endpointPolicy.Len()

		if err := endpointPolicy.Ready(); err != nil {
			logrus.WithField("podLabels", lbls).WithError(err).Error("failed to make endpoint policy ready")
		} else {
			endpointPolicy.Detach()
		}
	}

	return identityEntries
}

// derived from regeneratePolicy() at https://github.com/cilium/cilium/blob/b967f939236d2ed29a501acba5e2d764bfcba8c0/pkg/endpoint/policy.go#L218-L267
func toEndpointPolicy(repo *policy.Repository, endpointOwner policy.PolicyOwner, securityIdentity *identity.Identity) (*policy.EndpointPolicy, error) {
	var selectorPolicy policy.SelectorPolicy
	// force policy recompute
	skipPolicyRevision := 0
	selectorPolicy, policyRevision, err := repo.GetSelectorPolicy(securityIdentity, uint64(skipPolicyRevision), noOpStats{})
	_ = policyRevision
	if err != nil {
		// e.getLogger().WithError(err).Warning("Failed to calculate SelectorPolicy")
		return nil, err
	}

	// selectorPolicy is nil if skipRevision was matched.
	if selectorPolicy == nil {
		fmt.Printf("SelectorPolicy is nil\n")
		// e.getLogger().WithFields(logrus.Fields{
		// 	"policyRevision.next": e.nextPolicyRevision,
		// 	"policyRevision.repo": result.policyRevision,
		// 	"policyChanged":       e.nextPolicyRevision > e.policyRevision,
		// }).Debug("Skipping unnecessary endpoint policy recalculation")
		// datapathRegenCtxt.policyResult = result
		return nil, nil
	}

	// DistillPolicy converts a SelectorPolicy in to an EndpointPolicy
	endpointPolicy := selectorPolicy.DistillPolicy(endpointOwner, nil)
	return endpointPolicy, nil
}
