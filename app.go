package dot

import (
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/darian/dot/domain"
	"github.com/darian/dot/ui/projections"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// AppKey is the DoT application's identity key.
const AppKey = "dot-smart-home-0000-0000-0000-000000000001"

// App is a dogma.Application implementation for the DoT IoT smart home system.
type App struct {
	DeviceAggregate domain.DeviceHandler

	ReadDB             *sql.DB
	DeviceStateProjection projections.DeviceStateProjectionHandler
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("dot", AppKey)

	c.Routes(
		dogma.ViaAggregate(a.DeviceAggregate),

		dogma.ViaProjection(sqlprojection.New(a.ReadDB, sqlprojection.SQLiteDriver,
			&a.DeviceStateProjection)),
	)
}
