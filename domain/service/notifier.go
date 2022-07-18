package service

import "github.com/toledoom/test-pager/v2/domain/model"

type MailNotifier interface {
	Notifier
}

type SmsNotifier interface {
	Notifier
}

type Notifier interface {
	Notify(targets ...model.Target)
}

type CompositeNotifier struct {
	smsNotifier  SmsNotifier
	mailNotifier MailNotifier
}

func NewCompositeNotifier(smsNotifier SmsNotifier, mailNotifier MailNotifier) *CompositeNotifier {
	return &CompositeNotifier{
		smsNotifier:  smsNotifier,
		mailNotifier: mailNotifier,
	}
}

func (cn *CompositeNotifier) Notify(targets ...model.Target) {
	for _, target := range targets {
		switch target.(type) {
		case *model.EmailTarget:
			cn.mailNotifier.Notify(target)
		case *model.SmsTarget:
			cn.smsNotifier.Notify(target)
		}
	}
}
