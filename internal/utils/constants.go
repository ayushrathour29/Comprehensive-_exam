package utils

const (
	// Environment variable names
	EnvDatabaseURL = "DATABASE_URL"
	EnvPort        = "PORT"

	// Log messages
	LogNoEnvFileFound      = "No .env file found, relying on environment variables"
	LogDBConnectFail       = "Failed to connect to database"
	LogDBConnected         = "Connected to PostgreSQL database"
	LogDBClosed            = "Database connection closed"
	LogJobCreated          = "Job created"
	LogJobNotFound         = "Job not found"
	LogFailedToListJobs    = "Failed to list jobs"
	LogFailedToInsertJob   = "Failed to insert job"
	LogFailedToRetrieveJob = "Failed to retrieve job"
	LogFailedToScanJobRow  = "Failed to scan job row"
	LogRowIterationError   = "Row iteration error"

	// Error messages
	ErrInvalidPayload    = "Invalid payload"
	ErrFailedToCreateJob = "Failed to create job"
	ErrFailedToListJobs  = "Failed to list jobs"
)
