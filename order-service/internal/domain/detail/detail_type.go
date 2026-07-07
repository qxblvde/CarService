package detail

type DetailType string

const (
	Wheels       DetailType = "wheels"
	Transmission DetailType = "transmission"
	Engine       DetailType = "engine"
	Interior     DetailType = "interior"
	Steering     DetailType = "steering"
)

func (t DetailType) String() string {
	return string(t)
}

func (t DetailType) IsValid() bool {
	switch t {
	case Wheels, Transmission, Engine, Interior, Steering:
		return true
	default:
		return false
	}
}
