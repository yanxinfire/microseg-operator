package types

const (
	Ingress string = "Ingress"
	Egress         = "Egress"

	Segment   = "Segment"
	Resource  = "Resource"
	IPBlock   = "IPBlock"
	Isolation = "Isolation"

	Allow = "Allow"
	Deny  = "Deny"
	Log   = "Log"

	ProtocolUDP     = "UDP"
	ProtocolTCP     = "TCP"
	ProtocolICMP    = "ICMP"
	ProtocolICMPv6  = "ICMPv6"
	ProtocolSCTP    = "SCTP"
	ProtocolUDPLite = "UDPLite"
)
