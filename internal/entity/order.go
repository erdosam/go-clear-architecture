package entity

// OrderState represents the current state of a trash collection order
type OrderState string

const (
	StateCreated    OrderState = "created"     //Initial order creation
	StatePending    OrderState = "pending"     //Awaiting confirmation/payment
	StateConfirmed  OrderState = "confirmed"   //Order verified and scheduled
	StateInProgress OrderState = "in_progress" //Collection in process
	StateCompleted  OrderState = "completed"   //Collection finished successfully
	StateCancelled  OrderState = "cancelled"   //Order cancelled by customer

	StatePaymentPending   OrderState = "payment_pending"    //Waiting for payment processing
	StatePaymentFailed    OrderState = "payment_failed"     //When payment attempt unsuccessful
	StateScheduled        OrderState = "scheduled"          //When assigned a specific pickup date/time
	StateDriverAssigned   OrderState = "driver_assigned"    //When a specific collector is allocated
	StateDriverEnRoute    OrderState = "driver_en_route"    //Driver is heading to pickup location
	StateDriverAtLocation OrderState = "driver_at_location" //Driver is heading to pickup location
	StateDisputed         OrderState = "disputed"           //Customer raised issues with service
	StateRescheduled      OrderState = "rescheduled"        //Pickup date changed
	StateRefunded         OrderState = "refunded"           //Payment returned to customer

	StateActive  OrderState = "active"  //For ongoing service subscriptions
	StatePaused  OrderState = "paused"  //Temporary hold on recurring service
	StateExpired OrderState = "expired" //Order cancelled by system
	StateRenewed OrderState = "renewed" //Subscription extended by customer
)

type Order struct {
	Id        string `db:"id" json:"id"`
	LastState OrderState
	Timestamp
}

type OrderResponse struct {
}

// IsValidState checks if a given state is valid
func (o *Order) IsValidState() bool {
	validStates := map[OrderState]bool{
		StateCreated:          true,
		StatePending:          true,
		StateConfirmed:        true,
		StateInProgress:       true,
		StateCompleted:        true,
		StateCancelled:        true,
		StatePaymentPending:   true,
		StatePaymentFailed:    true,
		StateScheduled:        true,
		StateDriverAssigned:   true,
		StateDriverEnRoute:    true,
		StateDriverAtLocation: true,
		StateDisputed:         true,
		StateRescheduled:      true,
		StateRefunded:         true,
		StateActive:           true,
		StatePaused:           true,
		StateExpired:          true,
		StateRenewed:          true,
	}
	return validStates[o.LastState]
}
