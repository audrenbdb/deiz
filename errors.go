package deiz

const GenericError Error = "fail"
const ErrorClinicianDoesNotExist = "ce clinicien n'existe pas"
const ErrorUnauthorized Error = "unauthorized"
const ErrorStructValidation Error = "unable to validate struct"
const ErrorBookingSlotAlreadyFilled Error = "Opération incomplète, les créneaux n'étaient pas tous libres"

type Error string

func (e Error) Error() string {
	return string(e)
}
