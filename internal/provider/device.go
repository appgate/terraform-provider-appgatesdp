package provider

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
)

// defaultDeviceID derives a stable login device ID from the machine ID, so
// repeated runs on one host are seen as the same device; it falls back to a
// random UUID.
func defaultDeviceID() string {
	if id, err := machineid.ID(); err == nil {
		if u, err := uuid.Parse(id); err == nil {
			return u.String()
		}
	}
	return uuid.NewString()
}
