package billing

type repo interface {
	UnpaidBookingsGetter
	ClinicianBoundChecker
	BookingInvoiceCreater
	InvoicesCounter
	PaymentMethodsGetter
	PeriodInvoicesGetter
}

type mailer interface {
	BookingInvoiceMailer
	InvoicesSummaryMailer
}

type pdfer interface {
	BookingInvoicePDFCreater
	InvoicesSummaryPDFCreater
}

type Usecase struct {
	UnpaidBookingsGetter      UnpaidBookingsGetter
	ClinicianBoundChecker     ClinicianBoundChecker
	BookingInvoiceCreater     BookingInvoiceCreater
	BookingInvoiceMailer      BookingInvoiceMailer
	InvoicesCounter           InvoicesCounter
	BookingInvoicePDFCreater  BookingInvoicePDFCreater
	PaymentMethodsGetter      PaymentMethodsGetter
	PeriodInvoicesGetter      PeriodInvoicesGetter
	InvoicesSummaryPDFCreater InvoicesSummaryPDFCreater
	InvoicesSummaryMailer     InvoicesSummaryMailer
}

func NewUsecase(repo repo, mailer mailer, pdfer pdfer) *Usecase {
	return &Usecase{
		UnpaidBookingsGetter:      repo,
		ClinicianBoundChecker:     repo,
		BookingInvoiceCreater:     repo,
		BookingInvoiceMailer:      mailer,
		InvoicesSummaryMailer:     mailer,
		InvoicesCounter:           repo,
		BookingInvoicePDFCreater:  pdfer,
		InvoicesSummaryPDFCreater: pdfer,
		PaymentMethodsGetter:      repo,
		PeriodInvoicesGetter:      repo,
	}
}
