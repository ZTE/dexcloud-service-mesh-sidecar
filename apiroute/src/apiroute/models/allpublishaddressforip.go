package models

type AllPulishAddressForIP struct {
	IP              string `json:"ip"`
	Port            string `json:"port"`
	PublishProtocol string `json:"publish_protocol"`
}
