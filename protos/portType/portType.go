package portType

type PortType string

const (
	// HEALTH represents a health check port
	HEALTH PortType = "health"

	// API represents an API port
	API PortType = "api"

	// UNKNOWN represents an unknown port type
	UNKNOWN PortType = "unknown"
)

// String returns the string representation of the PortType
func (pt PortType) String() string {
	switch pt {
	case HEALTH:
		return "health"
	case API:
		return "api"
	default:
		return "unknown"
	}
}

// IsValid checks if the PortType is valid
func (pt PortType) IsValid() bool {
	return pt == HEALTH || pt == API
}

// ParsePortType parses a string into a PortType
func ParsePortType(s string) PortType {
	switch s {
	case "health":
		return HEALTH
	case "api":
		return API
	default:
		return UNKNOWN
	}
}

// CheckPortType checks if the given port type is valid
func (pt PortType) CheckPortType() bool {
	return pt.IsValid()
}
