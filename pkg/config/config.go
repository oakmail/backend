package config

// Config is the root config structure
type Config struct {
	LogLevel      string `default:"debug"`
	DefaultDomain string `default:"oakmail.io"`

	API    APIConfig
	Mailer MailerConfig
	Worker WorkerConfig

	Database *Database `default:"sqlite"`
	Postgres PostgresConfig
	SQLite   SQLiteConfig

	Filesystem *Filesystem `default:"flat"`
	Flat       FlatConfig
	Seaweed    SeaweedConfig

	Queue *Queue `default:"memory"`
	NSQ   NSQConfig
}

// APIConfig contains configuration data for the API module
type APIConfig struct {
	Enabled bool   `default:"false"`
	Address string `default:"0.0.0.0:8080"`
}

// MailerConfig contains configuration data for the mailer module
type MailerConfig struct {
	Enabled bool   `default:"false"`
	Address string `default:"0.0.0.0:8025"`
}

// WorkerConfig contains configuration data for the worker module
type WorkerConfig struct {
	Enabled bool `default:"false"`
}

// Database is a database type enum
type Database string

// String implements flag.Value
func (d *Database) String() string {
	return string(*d)
}

// Set implements flag.Value
func (d *Database) Set(value string) error {
	*d = Database(value)
	return nil
}

// Available databases
const (
	Postgres Database = "postgres"
	SQLite   Database = "sqlite"
)

// PostgresConfig contains all configuration data for a PostgreSQL connection
type PostgresConfig struct {
	ConnectionString string
}

// SQLiteConfig contains all configuration data for SQLite adapter setup
type SQLiteConfig struct {
	ConnectionString string `default:"./_runtime/database.db"`
}

// Filesystem is a filesystem type enum
type Filesystem string

// String implements flag.Value
func (f *Filesystem) String() string {
	return string(*f)
}

// Set implements flag.Value
func (f *Filesystem) Set(value string) error {
	*f = Filesystem(value)
	return nil
}

// Available filesystems
const (
	Flat    Filesystem = "flat"
	Seaweed Filesystem = "seaweed"
)

// FlatConfig contains all configuration data for the flat filesystem
type FlatConfig struct {
	Path string `default:"./.trtl/files"`
}

// SeaweedConfig contains all configuration data for the SeaweedFS client
type SeaweedConfig struct {
	MasterURL  string
	Collection string
}

// Queue is a queue type enum
type Queue string

// String implements flag.Value
func (q *Queue) String() string {
	return string(*q)
}

// Set implements flag.Value
func (q *Queue) Set(value string) error {
	*q = Queue(value)
	return nil
}

// Available queues
const (
	NSQ    Queue = "nsq"
	Memory Queue = "memory"
)

// NSQConfig contains all configuration data for the NSQ clients
type NSQConfig struct {
	NSQdAddresses    []string
	LookupdAddresses []string
}
