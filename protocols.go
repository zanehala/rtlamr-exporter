package main

type SCMPlus struct {
	ProtocolID   uint8  `json:"ProtocolID"`
	EndpointType uint8  `json:"EndpointType"`
	EndpointID   uint32 `json:"EndpointID"`
	Consumption  uint32 `json:"Consumption"`
	Tamper       uint16 `json:"Tamper"`
}

type SCM struct {
	ID          uint32 `json:"ID"`
	Type        uint8  `json:"Type"`
	TamperPhy   uint8  `json:"TamperPhy"`
	TamperEnc   uint8  `json:"TamperEnc"`
	Consumption uint32 `json:"Consumption"`
}

// Also used for NetIDM messages as the payload is roughly the same
type IDM struct {
	ERTSerialNumber uint32 `json:"ERTSerialNumber"`
	ERTType         uint8  `json:"ERTType"`
	// LastConsumptionCount is effectively equivalent to the Consumption count on SCM messages
	// It's a cumulative, monotonically increasing, consumption record
	LastConsumptionCount uint32 `json:"LastConsumptionCount"`
}

// Also used for R900BCD messages
type R900 struct {
	ID uint8 `json:"ID"`
	// ERT meter type
	Unkn1       uint8  `json:"Unkn1"`
	Consumption uint32 `json:"Consumption"`
}
