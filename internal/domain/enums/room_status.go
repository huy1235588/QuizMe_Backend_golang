package enums

// RoomStatus represents the status of a game room
type RoomStatus string

const (
	RoomStatusWaiting    RoomStatus = "WAITING"
	RoomStatusInProgress RoomStatus = "IN_PROGRESS"
	RoomStatusCompleted  RoomStatus = "COMPLETED"
	RoomStatusCancelled  RoomStatus = "CANCELLED"
)

// IsValid checks if the room status is valid
func (rs RoomStatus) IsValid() bool {
	switch rs {
	case RoomStatusWaiting, RoomStatusInProgress, RoomStatusCompleted, RoomStatusCancelled:
		return true
	}
	return false
}

// String returns the string representation
func (rs RoomStatus) String() string {
	return string(rs)
}
