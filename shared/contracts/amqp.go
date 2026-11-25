package contracts

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

type WSMessageType string

// Routing keys - using consistent event/command patterns
const (
	// Trip events (trip.event.*)
	TripEventCreated             WSMessageType = "trip.event.created"
	TripEventDriverAssigned      WSMessageType = "trip.event.driver_assigned"
	TripEventNoDriversFound      WSMessageType = "trip.event.no_drivers_found"
	TripEventDriverNotInterested WSMessageType = "trip.event.driver_not_interested"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest WSMessageType = "driver.cmd.trip_request"
	DriverCmdTripAccept  WSMessageType = "driver.cmd.trip_accept"
	DriverCmdTripDecline WSMessageType = "driver.cmd.trip_decline"
	DriverCmdLocation    WSMessageType = "driver.cmd.location"
	DriverCmdRegister    WSMessageType = "driver.cmd.register"

	// Payment events (payment.event.*)
	PaymentEventSessionCreated WSMessageType = "payment.event.session_created"
	PaymentEventSuccess        WSMessageType = "payment.event.success"
	PaymentEventFailed         WSMessageType = "payment.event.failed"
	PaymentEventCancelled      WSMessageType = "payment.event.cancelled"

	// Payment commands (payment.cmd.*)
	PaymentCmdCreateSession WSMessageType = "payment.cmd.create_session"
)
