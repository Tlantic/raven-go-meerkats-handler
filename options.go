package raven_meerkats

import (
	"github.com/Tlantic/meerkats"
	"github.com/getsentry/raven-go"
)

//noinspection GoUnusedExportedFunction
func Client(cli *raven.Client) meerkats.HandlerReceiver {
	return meerkats.HandlerReceiver(func(h meerkats.Handler) {
		h.(*RavenHandler).Client = cli
	})
}
