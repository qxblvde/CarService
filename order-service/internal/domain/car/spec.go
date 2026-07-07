package car

type BodyType string

const (
	BodyTypeSedan     BodyType = "sedan"
	BodyTypeWagon     BodyType = "wagon"
	BodyTypeCoupe     BodyType = "coupe"
	BodyTypeHatchback BodyType = "hatchback"
	BodyTypeSUV       BodyType = "suv"
	BodyTypeCrossover BodyType = "crossover"
	BodyTypeMinivan   BodyType = "minivan"
	BodyTypePickup    BodyType = "pickup"
)

type FuelType string

const (
	FuelTypePetrol   FuelType = "petrol"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeElectric FuelType = "electric"
	FuelTypeHybrid   FuelType = "hybrid"
)

type TransmissionType string

const (
	TransmissionManual    TransmissionType = "manual"
	TransmissionAutomatic TransmissionType = "automatic"
)

type DriveType string

const (
	DriveFront    DriveType = "front"
	DriveRear     DriveType = "rear"
	DriveAllWheel DriveType = "all_wheel"
)

func (t BodyType) String() string {
	return string(t)
}

func (t BodyType) IsValid() bool {
	switch t {
	case BodyTypeSedan, BodyTypeWagon, BodyTypeCoupe, BodyTypeHatchback, BodyTypeSUV, BodyTypeCrossover, BodyTypeMinivan, BodyTypePickup:
		return true
	default:
		return false
	}
}

func (t FuelType) String() string {
	return string(t)
}

func (t FuelType) IsValid() bool {
	switch t {
	case FuelTypePetrol, FuelTypeDiesel, FuelTypeElectric, FuelTypeHybrid:
		return true
	default:
		return false
	}
}

func (t TransmissionType) String() string {
	return string(t)
}

func (t TransmissionType) IsValid() bool {
	switch t {
	case TransmissionManual, TransmissionAutomatic:
		return true
	default:
		return false
	}
}

func (t DriveType) String() string {
	return string(t)
}

func (t DriveType) IsValid() bool {
	switch t {
	case DriveFront, DriveRear, DriveAllWheel:
		return true
	default:
		return false
	}
}
