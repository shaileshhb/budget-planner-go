package config

// EnvKey are the type of constants that will be used as keys while reading the .env file
type EnvKey string

// Constants
const (
	// For DB
	DBName         EnvKey = "DB_NAME"
	DBUser         EnvKey = "DB_USER"
	DBPass         EnvKey = "DB_PASSWORD"
	DBPort         EnvKey = "DB_PORT"
	DBHost         EnvKey = "DB_HOST"
	ENCRYPTION_KEY EnvKey = "ENCRYPTION_KEY"

	// For Server
	PORT             EnvKey = "PORT"
	JWTKey           EnvKey = "JWT_KEY"
	HTTPWriteTimeout EnvKey = "HTTP_WRITE_TIMEOUT"
	HTTPReadTimeout  EnvKey = "HTTP_READ_TIMEOUT"
	HTTPIdleTimeout  EnvKey = "HTTP_IDLE_TIMEOUT"

	// For emails
	SMTPHost   string = "smtp.gmail.com"
	SMTPPort   string = "587"
	SenderMail EnvKey = "SENDER_MAIL"
	SenderPass EnvKey = "SENDER_PASSWORD"
)
