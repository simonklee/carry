package tsdb

import (
	"time"

	"github.com/influxdb/influxdb/toml"
)

const (
	// DefaultRetentionAutoCreate is the default for auto-creating retention policies
	DefaultRetentionAutoCreate = true

	// DefaultRetentionCheckEnabled is the default for checking for retention policy enforcement
	DefaultRetentionCheckEnabled = true

	// DefaultRetentionCreatePeriod represents how often the server will check to see if new
	// shard groups need to be created in advance for writing
	DefaultRetentionCreatePeriod = 45 * time.Minute

	// DefaultRetentionCheckPeriod is the period of time between retention policy checks are run
	DefaultRetentionCheckPeriod = 10 * time.Minute
)

type Config struct {
	Dir                   string        `toml:"dir"`
	RetentionAutoCreate   bool          `toml:"retention-auto-create"`
	RetentionCheckEnabled bool          `toml:"retention-check-enabled"`
	RetentionCheckPeriod  toml.Duration `toml:"retention-check-period"`
	RetentionCreatePeriod toml.Duration `toml:"retention-create-period"`
}

func NewConfig() Config {
	return Config{
		RetentionAutoCreate:   DefaultRetentionAutoCreate,
		RetentionCheckEnabled: DefaultRetentionCheckEnabled,
		RetentionCheckPeriod:  toml.Duration(DefaultRetentionCheckPeriod),
		RetentionCreatePeriod: toml.Duration(DefaultRetentionCreatePeriod),
	}
}

// ShardGroupPreCreateCheckPeriod returns the check interval to pre-create shard groups.
// If it was not defined in the config, it defaults to DefaultShardGroupPreCreatePeriod
func (c *Config) ShardGroupPreCreateCheckPeriod() time.Duration {
	if c.RetentionCreatePeriod != 0 {
		return time.Duration(c.RetentionCreatePeriod)
	}
	return DefaultRetentionCreatePeriod
}
