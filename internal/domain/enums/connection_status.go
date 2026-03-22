package enums

// ConnectionStatus represents the connection status of a participant
type ConnectionStatus string

const (
	ConnectionStatusActive       ConnectionStatus = "ACTIVE"
	ConnectionStatusDisconnected ConnectionStatus = "DISCONNECTED"
	ConnectionStatusTimedOut     ConnectionStatus = "TIMED_OUT"
)

func (s ConnectionStatus) String() string {
	return string(s)
}
