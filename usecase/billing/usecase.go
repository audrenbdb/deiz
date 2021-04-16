package billing

type repo interface {
	UnpaidBookingsGetter
	ClinicianBoundChecker
	BookingInvoiceCreater
	InvoicesCounter
	PaymentMethodsGetter
	PeriodInvoicesGetter
	ClinicianStripeSecretKeyGetter
	BookingInvoiceCanceler
}

type mailer interface {
	BookingInvoiceMailer
	InvoicesSummaryMailer
}

type pdfer interface {
	BookingInvoicePDFCreater
	InvoicesSummaryPDFCreater
}

type crypter interface {
	BytesDecrypter
}

type striper interface {
	StripePaymentSessionCreater
}

type Usecase struct {
	UnpaidBookingsGetter           UnpaidBookingsGetter
	ClinicianBoundChecker          ClinicianBoundChecker
	BookingInvoiceCreater          BookingInvoiceCreater
	BookingInvoiceMailer           BookingInvoiceMailer
	InvoicesCounter                InvoicesCounter
	BookingInvoicePDFCreater       BookingInvoicePDFCreater
	PaymentMethodsGetter           PaymentMethodsGetter
	PeriodInvoicesGetter           PeriodInvoicesGetter
	InvoicesSummaryPDFCreater      InvoicesSummaryPDFCreater
	InvoicesSummaryMailer          InvoicesSummaryMailer
	ClinicianStripeSecretKeyGetter ClinicianStripeSecretKeyGetter
	BytesDecrypter                 BytesDecrypter
	StripePaymentSessionCreater    StripePaymentSessionCreater
	BookingInvoiceCanceler         BookingInvoiceCanceler
}

func NewUsecase(repo repo, mailer mailer, pdfer pdfer, crypter crypter, striper striper) *Usecase {
	return &Usecase{
		UnpaidBookingsGetter:           repo,
		ClinicianBoundChecker:          repo,
		BookingInvoiceCreater:          repo,
		BookingInvoiceMailer:           mailer,
		InvoicesSummaryMailer:          mailer,
		InvoicesCounter:                repo,
		BookingInvoicePDFCreater:       pdfer,
		InvoicesSummaryPDFCreater:      pdfer,
		PaymentMethodsGetter:           repo,
		PeriodInvoicesGetter:           repo,
		ClinicianStripeSecretKeyGetter: repo,
		BytesDecrypter:                 crypter,
		StripePaymentSessionCreater:    striper,
		BookingInvoiceCanceler:         repo,
	}
}
