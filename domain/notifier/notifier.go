package notifier

import (
	"github.com/toledoom/test-pager/v2/domain/escalationpolicy"
)

type Mail interface {
	Notifier
}

type Sms interface {
	Notifier
}

type Notifier interface {
	Notify(targets ...escalationpolicy.Target)
}

type Composite struct {
	sms  Sms
	mail Mail
}

func NewComposite(smsNotifier Sms, mailNotifier Mail) *Composite {
	return &Composite{
		sms:  smsNotifier,
		mail: mailNotifier,
	}
}

func (cn *Composite) Notify(targets ...escalationpolicy.Target) {
	for _, target := range targets {
		switch target.(type) {
		case *escalationpolicy.EmailTarget:
			cn.mail.Notify(target)
		case *escalationpolicy.SmsTarget:
			cn.sms.Notify(target)
		}
	}
}
