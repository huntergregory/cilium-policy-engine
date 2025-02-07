package main

import (
	"github.com/huntergregory/cilium-policy-engine/cmd/inputs"
	"github.com/huntergregory/cilium-policy-engine/pkg/policyengine/identities"
	"github.com/huntergregory/cilium-policy-engine/pkg/policyengine/policies"
	"github.com/sirupsen/logrus"

	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"

	"github.com/cilium/cilium/pkg/option"
	"github.com/cilium/cilium/pkg/policy"
	"github.com/cilium/cilium/pkg/policy/api"
	// "github.com/cilium/cilium/pkg/logging"
)

// NOTE: two flaws currently:
// 1. Any policy add will readd the policy due to limitations in policy.Repository.MustAddList(). A workaround is to temporarily modify vendor/ to reset p.nextID to 0 for each addListLocked() call.
// 2. There is still some leak causing eventual sigkills if you keep (re)creating netpols and/or keep generating endpoint policies.

func main() {
	// logging.DefaultLogger.SetLevel(logrus.DebugLevel)
	policy.SetPolicyEnabled(option.DefaultEnforcement)
	repo := policy.NewPolicyRepository(nil, nil, nil, nil, api.NewPolicyMetricsNoop())
	lblsToIdentity := identities.CreateIdentities(repo, inputs.Pods)
	policies.AddPolicies(repo, inputs.NetPols)
	identityEntries := policies.CountEntries(repo, lblsToIdentity)

	maxEntries := 0
	maxLbls := ""
	for lbls, entries := range identityEntries {
		logrus.WithFields(logrus.Fields{
			"podLabels": lbls,
			"entries":   entries,
		}).Debug("calculated policy map entries for pod identity")

		if entries > maxEntries {
			maxEntries = entries
			maxLbls = lbls
		}
	}

	logrus.WithFields(logrus.Fields{
		"maxPodLabels": maxLbls,
		"maxEntries":   maxEntries,
	}).Info("calculated max policy map entries for all pod identities")
}
