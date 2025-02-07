package main

import (
	"fmt"
	"time"

	"github.com/huntergregory/cilium-policy-engine/cmd/inputs"
	"github.com/huntergregory/cilium-policy-engine/pkg/policyengine/identities"
	"github.com/huntergregory/cilium-policy-engine/pkg/policyengine/policies"
	"github.com/sirupsen/logrus"

	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"

	slim_networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/networking/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/util/intstr"
	"github.com/cilium/cilium/pkg/option"
	"github.com/cilium/cilium/pkg/policy"
	"github.com/cilium/cilium/pkg/policy/api"
	// "github.com/cilium/cilium/pkg/logging"
)

// NOTE: two flaws currently:
// 1. Any policy add will readd the policy due to limitations in policy.Repository.MustAddList(). A workaround is to temporarily modify vendor/ to reset p.nextID to 0 for each addListLocked() call.
// 2. There is still some GC leak causing eventual sigkills if you add several netpols and/or keep generating endpoint policies.
//
// The write/read to file code commented out below was a hack to deal with these sigkills.

func main() {
	// logging.DefaultLogger.SetLevel(logrus.DebugLevel)
	policy.SetPolicyEnabled(option.DefaultEnforcement)

	repo := policy.NewPolicyRepository(nil, nil, nil, nil, api.NewPolicyMetricsNoop())
	lblsToIdentity := identities.CreateIdentities(repo, inputs.Pods)

	// // read starting port from file
	// startingPort := 0
	// // use os package
	// file, err := os.Open("startingPort.txt")
	// if err != nil {
	// 	logrus.WithError(err).Panic("failed to open file")
	// }
	// _, err = fmt.Fscanf(file, "%d", &startingPort)
	// if err != nil {
	// 	logrus.WithError(err).Panic("failed to read starting port from file")
	// }
	// err = file.Close()
	// if err != nil {
	// 	logrus.WithError(err).Panic("failed to close file")
	// }

	// fmt.Println("Starting port is:", startingPort)

	maxEntries := 0
	maxLbls := ""
	maxPort := 0
	maxEndPort := 0
	// for port := 65535; port >= 1; port-- {
	for port := 65535; port >= 1; port-- {
		if port%20 == 0 {
			// logrus.WithFields(logrus.Fields{
			// 	"date":         time.Now(),
			// 	"port":         port,
			// 	"maxPodLabels": maxLbls,
			// 	"maxEntries":   maxEntries,
			// 	"maxPort":      maxPort,
			// 	"maxEndPort":   maxEndPort,
			// }).Info("processing port")

			// Print the same thing
			fmt.Printf("date: %v, port: %v, maxPodLabels: %v, maxEntries: %v, maxPort: %v, maxEndPort: %v\n", time.Now(), port, maxLbls, maxEntries, maxPort, maxEndPort)

			// // overwrite current port to file
			// file, err := os.Create("startingPort.txt")
			// if err != nil {
			// 	logrus.WithError(err).Panic("failed to open file")
			// }
			// _, err = fmt.Fprintf(file, "%d", port)
			// if err != nil {
			// 	logrus.WithError(err).Panic("failed to write starting port to file")
			// }
			// err = file.Close()
			// if err != nil {
			// 	logrus.WithError(err).Panic("failed to close file")
			// }
		}

		// for endPort := port + 1; endPort <= 65535; endPort++ {
		endPort := 65534
		if endPort <= port {
			continue
		}
		policies.AddPolicies(repo, []*slim_networkingv1.NetworkPolicy{netPol(port, endPort)})
		identityEntries := policies.CountEntries(repo, lblsToIdentity)

		for lbls, entries := range identityEntries {
			// logrus.WithFields(logrus.Fields{
			// 	"podLabels": lbls,
			// 	"entries":   entries,
			// }).Trace("calculated policy map entries for pod identity")

			if entries > maxEntries {
				maxEntries = entries
				maxLbls = lbls
				maxPort = port
				maxEndPort = endPort
			}
		}

		// logrus.WithFields(logrus.Fields{
		// 	"maxPodLabels":   maxLbls,
		// 	"maxEntries":     maxEntries,
		// 	"maxPort":        maxPort,
		// 	"maxEndPort":     maxEndPort,
		// 	"currentPort":    port,
		// 	"currentEndPort": endPort,
		// }).Debug("calculated max policy map entries for all pod identities")

		// if err := repo.GetSelectorCache().GetVersionHandle().Close(); err != nil {
		// 	logrus.WithError(err).Error("failed to close version handle")
		// 	return
		// }
		// }
	}

	logrus.WithFields(logrus.Fields{
		"maxPodLabels": maxLbls,
		"maxEntries":   maxEntries,
		"maxPort":      maxPort,
		"maxEndPort":   maxEndPort,
	}).Info("calculated max policy map entries for all pod identities and port ranges")
}

func netPol(port, endPort int) *slim_networkingv1.NetworkPolicy {
	end32 := int32(endPort)
	end32Pointer := &end32
	return &slim_networkingv1.NetworkPolicy{
		ObjectMeta: slim_metav1.ObjectMeta{
			Name:      "test-egress",
			Namespace: "default",
		},
		Spec: slim_networkingv1.NetworkPolicySpec{
			PodSelector: slim_metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test",
				},
			},
			Egress: []slim_networkingv1.NetworkPolicyEgressRule{
				{
					To: []slim_networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &slim_metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "test",
								},
							},
						},
					},
					Ports: []slim_networkingv1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: int32(port),
							},
							EndPort: end32Pointer,
						},
					},
				},
			},
			PolicyTypes: []slim_networkingv1.PolicyType{
				slim_networkingv1.PolicyTypeEgress,
			},
		},
	}
}
