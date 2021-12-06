# prisma-sdk-v2

```golang
	egressRule := prisma_types.NewRule().WithTrafficActionAllow().AddTCPProtocolPort(53).AddUDPProtocolPort(53).
		AddObject("@org:cloudaccount=prisma-microseg-field-azure", "@org:group=prisma-aks-microseg", "@org:tenant=806775361903163392", "externalnetwork:name=Kube DNS")

	netRuleSet := prisma_types.NewNetworkrulesetpolicy("Kube DNS").
		WithDescription("auto-generated cloud operator policy").
		AddOutgoingRule(egressRule).
		AddSubject("@org:cloudaccount=prisma-microseg-field-azure", "@org:group=prisma-aks-microseg", "@org:tenant=806775361903163392").
		WithPropagate(true).WithProtected(true)
```


