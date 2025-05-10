package models

type CertstreamMessage struct {
	MessageType string `json:"message_type"`
	Data        struct {
		UpdateType string `json:"update_type"`
		LeafCert   struct {
			NotBefore  float64  `json:"not_before"`
			NotAfter   float64  `json:"not_after"`
			AllDomains []string `json:"all_domains"`
		} `json:"leaf_cert"`
	} `json:"data"`
}
