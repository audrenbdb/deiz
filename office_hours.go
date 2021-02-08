package deiz

import "context"

type OfficeHours struct {
	ID            int     `json:"id" validator:"required"`
	StartMn       int     `json:"startMn" validator:"required"`
	EndMn         int     `json:"endMn" validator:"required"`
	WeekDay       int     `json:"weekDay" validator:"required"`
	Address       Address `json:"address" validator:"required"`
	RemoteAllowed bool    `json:"remoteAllowed"`
}

//repo interfaces
type (
	officeHoursGetter interface {
		GetOfficeHours(ctx context.Context, clinicianID int) ([]OfficeHours, error)
	}
	officeHoursAdder interface {
		AddOfficeHours(ctx context.Context, h *OfficeHours, clinicianID int) error
	}
	officeHoursRemover interface {
		RemoveOfficeHours(ctx context.Context, h *OfficeHours, clinicianID int) error
	}
)

type (
	AddOfficeHours    func(ctx context.Context, h *OfficeHours, clinicianID int) error
	RemoveOfficeHours func(ctx context.Context, h *OfficeHours, clinicianID int) error
)

func addOfficeHoursFunc(adder officeHoursAdder) AddOfficeHours {
	return func(ctx context.Context, h *OfficeHours, clinicianID int) error {
		return adder.AddOfficeHours(ctx, h, clinicianID)
	}
}

func removeOfficeHoursFunc(remover officeHoursRemover) RemoveOfficeHours {
	return func(ctx context.Context, h *OfficeHours, clinicianID int) error {
		return remover.RemoveOfficeHours(ctx, h, clinicianID)
	}
}
