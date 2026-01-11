package conntrack

var ctStatusFlags = map[uint32]string{
	0x0001: "EXPECTED",      // Connexion créée car attendue (ex: flux data d’un FTP, SIP, etc.)
	0x0002: "SEEN_REPLY",    // Une réponse a été vue (le flux est bidirectionnel)
	0x0004: "ASSURED",       // Connexion considérée comme établie et fiable
	0x0008: "CONFIRMED",     // Entrée confirmée dans la table conntrack
	0x0010: "SRC_NAT",       // Source NAT appliqué (SNAT actif)
	0x0020: "DST_NAT",       // Destination NAT appliqué (DNAT actif)
	0x0040: "SEQ_ADJUST",    // Ajustement des numéros de séquence TCP (NAT ou helper)
	0x0080: "SRC_NAT_DONE",  // SNAT déjà effectué (état final appliqué)
	0x0100: "DST_NAT_DONE",  // DNAT déjà effectué (état final appliqué)
	0x0200: "DYING",         // Connexion en cours de suppression (timeout ou teardown)
	0x0400: "FIXED_TIMEOUT", // Timeout fixé manuellement (helper ou règle spécifique)
	0x0800: "TEMPLATE",      // Entrée modèle servant à créer d’autres connexions
}

type IPTupple struct {
	SrcIP   string
	DstIP   string
	SrcPort int
	DstPort int
}

type ConnectionTracked struct {
	Origin IPTupple
	Return IPTupple
	Status []string
}
