package provider

import (
	"errors"
	"os"
	"reflect"

	"github.com/microapis/messages-api"
)

// EmailProvider ...
type EmailProvider messages.Provider

// Load ...
func Load(p interface{}) error {
	var e reflect.Value

	switch p.(EmailProvider).Name {
	case "ses":
		e = reflect.ValueOf(p.(SESEmailProvider).Params.(ParamsSES)).Elem()
	case "mandrill":
		e = reflect.ValueOf(p.(MandrillEmailProvider).Params.(ParamsMandrill)).Elem()
	case "sendgrid":
		e = reflect.ValueOf(p.(SengridEmailProvider).Params.(ParamsSendGrid)).Elem()
	}

	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		tag, ok := e.Type().Field(i).Tag.Lookup("env")
		if !ok {
			return errors.New("key name: " + name + " is not valid")
		}

		value := os.Getenv("PROVIDER_" + tag)
		if value == "" {
			return errors.New("PROVIDER_" + tag + " env value not defined")
		}

		switch p.(EmailProvider).Name {
		case "ses":
			params, ok := p.(SESEmailProvider).Params.(ParamsSES)
			if !ok {
				return errors.New("key name: " + name + " could not cast ParamsSES")
			}
			params.SetParam(name, value)
		case "mandrill":
			params, ok := p.(MandrillEmailProvider).Params.(ParamsMandrill)
			if !ok {
				return errors.New("key name: " + name + " could not cast ParamsMandrill")
			}
			params.SetParam(name, value)
		case "sendgrid":
			params, ok := p.(SengridEmailProvider).Params.(ParamsSendGrid)
			if !ok {
				return errors.New("key name: " + name + " could not cast ParamsSendGrid")
			}
			params.SetParam(name, value)
		}
	}

	return nil
}
