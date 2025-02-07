package inputs

import (

	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"

	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	slim_networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/networking/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/util/intstr"
	// "github.com/cilium/cilium/pkg/logging"
)

var (
	udp = slim_corev1.ProtocolUDP
	tcp = slim_corev1.ProtocolTCP

	Pods = []*slim_corev1.Pod{
		{
			ObjectMeta: slim_metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		{
			ObjectMeta: slim_metav1.ObjectMeta{
				Name:      "test2",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test2",
				},
			},
			Spec: slim_corev1.PodSpec{
				Containers: []slim_corev1.Container{
					{
						Name: "cont1",
						// Ports: []slim_corev1.ContainerPort{
						// 	{
						// 		Name:          "http",
						// 		ContainerPort: 80,
						// 	},
						// },
					},
				},
			},
		},
	}

	NetPols = []*slim_networkingv1.NetworkPolicy{
		{
			ObjectMeta: slim_metav1.ObjectMeta{
				Name:      "test-egress-cidr",
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
									IntVal: 1,
								},
								EndPort: intPointer(65534),
								// EndPort: intPointer(32767 - 1),
								Protocol: &tcp,
							},
						},
					},
				},
				PolicyTypes: []slim_networkingv1.PolicyType{
					slim_networkingv1.PolicyTypeEgress,
				},
			},
		},
		// {
		// 	ObjectMeta: slim_metav1.ObjectMeta{
		// 		Name:      "test",
		// 		Namespace: "default",
		// 	},
		// 	Spec: slim_networkingv1.NetworkPolicySpec{
		// 		PodSelector: slim_metav1.LabelSelector{
		// 			MatchLabels: map[string]string{
		// 				"app": "test",
		// 			},
		// 		},
		// 		Ingress: []slim_networkingv1.NetworkPolicyIngressRule{
		// 			{
		// 				From: []slim_networkingv1.NetworkPolicyPeer{
		// 					{
		// 						PodSelector: &slim_metav1.LabelSelector{
		// 							MatchLabels: map[string]string{
		// 								"app": "test",
		// 							},
		// 						},
		// 					},
		// 				},
		// 				// Ports: []slim_networkingv1.NetworkPolicyPort{
		// 				// 	{
		// 				// 		Port: &intstr.IntOrString{
		// 				// 			Type:   intstr.Int,
		// 				// 			IntVal: 86,
		// 				// 		},
		// 				// 	},
		// 				// },
		// 			},
		// 		},
		// 		PolicyTypes: []slim_networkingv1.PolicyType{
		// 			slim_networkingv1.PolicyTypeIngress,
		// 		},
		// 	},
		// },
		// {
		// 	ObjectMeta: slim_metav1.ObjectMeta{
		// 		Name:      "test-duplicate",
		// 		Namespace: "default",
		// 	},
		// 	Spec: slim_networkingv1.NetworkPolicySpec{
		// 		PodSelector: slim_metav1.LabelSelector{
		// 			MatchLabels: map[string]string{
		// 				"app": "test",
		// 			},
		// 		},
		// 		Ingress: []slim_networkingv1.NetworkPolicyIngressRule{
		// 			{
		// 				From: []slim_networkingv1.NetworkPolicyPeer{
		// 					{
		// 						PodSelector: &slim_metav1.LabelSelector{
		// 							MatchLabels: map[string]string{
		// 								"app": "test",
		// 							},
		// 						},
		// 					},
		// 				},
		// 				Ports: []slim_networkingv1.NetworkPolicyPort{
		// 					{
		// 						Port: &intstr.IntOrString{
		// 							Type:   intstr.Int,
		// 							IntVal: 88,
		// 						},
		// 					},
		// 					// {
		// 					// 	Port: &intstr.IntOrString{
		// 					// 		Type:   intstr.Int,
		// 					// 		IntVal: 85,
		// 					// 	},
		// 					// 	// Protocol: &tcp,
		// 					// },
		// 					// {
		// 					// 	Port: &intstr.IntOrString{
		// 					// 		Type:   intstr.Int,
		// 					// 		IntVal: 89,
		// 					// 	},
		// 					// 	Protocol: &udp,
		// 					// },
		// 				},
		// 			},
		// 		},
		// 		PolicyTypes: []slim_networkingv1.PolicyType{
		// 			slim_networkingv1.PolicyTypeIngress,
		// 		},
		// 	},
		// },
		// {
		// 	ObjectMeta: slim_metav1.ObjectMeta{
		// 		Name:      "test2",
		// 		Namespace: "default",
		// 	},
		// 	Spec: slim_networkingv1.NetworkPolicySpec{
		// 		PodSelector: slim_metav1.LabelSelector{
		// 			MatchLabels: map[string]string{
		// 				"app": "test",
		// 			},
		// 		},
		// 		Ingress: []slim_networkingv1.NetworkPolicyIngressRule{
		// 			{
		// 				From: []slim_networkingv1.NetworkPolicyPeer{
		// 					{
		// 						PodSelector: &slim_metav1.LabelSelector{
		// 							MatchLabels: map[string]string{
		// 								"app": "test2",
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 		PolicyTypes: []slim_networkingv1.PolicyType{
		// 			slim_networkingv1.PolicyTypeIngress,
		// 		},
		// 	},
		// },
		// {
		// 	ObjectMeta: slim_metav1.ObjectMeta{
		// 		Name:      "test3-port",
		// 		Namespace: "default",
		// 	},
		// 	Spec: slim_networkingv1.NetworkPolicySpec{
		// 		PodSelector: slim_metav1.LabelSelector{
		// 			MatchLabels: map[string]string{
		// 				"app": "test",
		// 			},
		// 		},
		// 		Ingress: []slim_networkingv1.NetworkPolicyIngressRule{
		// 			{
		// 				From: []slim_networkingv1.NetworkPolicyPeer{
		// 					{
		// 						PodSelector: &slim_metav1.LabelSelector{
		// 							MatchLabels: map[string]string{
		// 								"app": "test3",
		// 							},
		// 						},
		// 					},
		// 				},
		// 				Ports: []slim_networkingv1.NetworkPolicyPort{
		// 					{
		// 						Port: &intstr.IntOrString{
		// 							Type:   intstr.Int,
		// 							IntVal: 82,
		// 						},
		// 					},
		// 					{
		// 						Port: &intstr.IntOrString{
		// 							Type:   intstr.Int,
		// 							IntVal: 85,
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 		PolicyTypes: []slim_networkingv1.PolicyType{
		// 			slim_networkingv1.PolicyTypeIngress,
		// 		},
		// 	},
		// },
		// {
		// 	ObjectMeta: slim_metav1.ObjectMeta{
		// 		Name:      "test-egress",
		// 		Namespace: "default",
		// 	},
		// 	Spec: slim_networkingv1.NetworkPolicySpec{
		// 		PodSelector: slim_metav1.LabelSelector{
		// 			MatchLabels: map[string]string{
		// 				"app": "test",
		// 			},
		// 		},
		// 		Egress: []slim_networkingv1.NetworkPolicyEgressRule{
		// 			{
		// 				To: []slim_networkingv1.NetworkPolicyPeer{
		// 					{
		// 						PodSelector: &slim_metav1.LabelSelector{
		// 							MatchLabels: map[string]string{
		// 								"app": "test",
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 		PolicyTypes: []slim_networkingv1.PolicyType{
		// 			slim_networkingv1.PolicyTypeEgress,
		// 		},
		// 	},
		// },
		// {
		// 	ObjectMeta: slim_metav1.ObjectMeta{
		// 		Name:      "test-egress-cidr",
		// 		Namespace: "default",
		// 	},
		// 	Spec: slim_networkingv1.NetworkPolicySpec{
		// 		PodSelector: slim_metav1.LabelSelector{
		// 			MatchLabels: map[string]string{
		// 				"app": "test",
		// 			},
		// 		},
		// 		Egress: []slim_networkingv1.NetworkPolicyEgressRule{
		// 			{
		// 				To: []slim_networkingv1.NetworkPolicyPeer{
		// 					{
		// 						IPBlock: &slim_networkingv1.IPBlock{
		// 							CIDR: "10.0.0.0/16",
		// 						},
		// 					},
		// 					{
		// 						IPBlock: &slim_networkingv1.IPBlock{
		// 							CIDR: "11.0.0.0/24",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 		PolicyTypes: []slim_networkingv1.PolicyType{
		// 			slim_networkingv1.PolicyTypeEgress,
		// 		},
		// 	},
		// },
	}
)

func intPointer(i int32) *int32 {
	return &i
}
