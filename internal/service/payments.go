package service

import (
	"context"
	"github.com/cookienyancloud/back/pkg/payment"
)

const (
	redirectURLTmpl = "https://%s/"
)

type PaymentsService struct {
	ordersService   Orders
	fondyCallbackURL string
}

func (s *PaymentsService) GeneratePaymentLink(ctx context.Context, orderId ) (string, error) {
	order, err := s.ordersService.GetById(ctx, orderId)
	if err != nil {
		return "", err
	}

	offer, err := s.offersService.GetById(ctx, order.Offer.ID)
	if err != nil {
		return "", err
	}

	if !offer.PaymentMethod.UsesProvider {
		return "", domain.ErrPaymentProviderNotUsed
	}

	paymentInput := payment.GeneratePaymentLinkInput{
		OrderId:   orderId.Hex(),
		Amount:    order.Amount,
		Currency:  offer.Price.Currency,
		OrderDesc: offer.Description, // TODO proper order description
	}

	switch offer.PaymentMethod.Provider {
	case domain.PaymentProviderFondy:
		return s.generateFondyPaymentLink(ctx, offer.SchoolID, paymentInput)
	default:
		return "", domain.ErrUnknownPaymentProvider
	}
}
