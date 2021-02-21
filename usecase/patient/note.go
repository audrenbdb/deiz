package patient

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	NotesGetter interface {
		GetPatientNotes(ctx context.Context, patientID int) ([]deiz.PatientNote, error)
	}
	NoteCreater interface {
		CreatePatientNote(ctx context.Context, n *deiz.PatientNote, patientID int) error
	}
	NoteDeleter interface {
		DeletePatientNote(ctx context.Context, noteID int, patientID int) error
	}
)

func (u *Usecase) GetPatientNotes(ctx context.Context, patientID int, clinicianID int) ([]deiz.PatientNote, error) {
	tied, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &deiz.Patient{ID: patientID}, clinicianID)
	if err != nil {
		return nil, err
	}
	if !tied {
		return nil, deiz.ErrorUnauthorized
	}
	return u.NotesGetter.GetPatientNotes(ctx, patientID)
}

func (u *Usecase) AddPatientNote(ctx context.Context, n *deiz.PatientNote, patientID int, clinicianID int) error {
	tied, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &deiz.Patient{ID: patientID}, clinicianID)
	if err != nil {
		return err
	}
	if !tied {
		return deiz.ErrorUnauthorized
	}
	return u.NoteCreater.CreatePatientNote(ctx, n, patientID)
}

func (u *Usecase) RemovePatientNote(ctx context.Context, noteID int, patientID int, clinicianID int) error {
	tied, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &deiz.Patient{ID: patientID}, clinicianID)
	if err != nil {
		return err
	}
	if !tied {
		return deiz.ErrorUnauthorized
	}
	return u.NoteDeleter.DeletePatientNote(ctx, noteID, patientID)
}
