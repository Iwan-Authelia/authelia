package storage

import (
	"database/sql"
	"fmt"

	"github.com/clems4ever/authelia/internal/logging"
	_ "github.com/mattn/go-sqlite3" // Load the SQLite Driver used in the connection string.
)

// SQLiteProvider is a sqlite3 provider
type SQLiteProvider struct {
	SQLProvider
}

// NewSQLiteProvider construct a sqlite provider.
func NewSQLiteProvider(path string) *SQLiteProvider {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		logging.Logger().Fatalf("Unable to create SQLite database %s: %s", path, err)
	}

	provider := SQLiteProvider{
		SQLProvider{
			sqlGetPreferencesByUsername:     fmt.Sprintf("SELECT second_factor_method FROM %s WHERE username=?", preferencesTableName),
			sqlUpsertSecondFactorPreference: fmt.Sprintf("REPLACE INTO %s (username, second_factor_method) VALUES (?, ?)", preferencesTableName),

			sqlTestIdentityVerificationTokenExistence: fmt.Sprintf("SELECT EXISTS (SELECT * FROM %s WHERE token=?)", identityVerificationTokensTableName),
			sqlInsertIdentityVerificationToken:        fmt.Sprintf("INSERT INTO %s (token) VALUES (?)", identityVerificationTokensTableName),
			sqlDeleteIdentityVerificationToken:        fmt.Sprintf("DELETE FROM %s WHERE token=?", identityVerificationTokensTableName),

			sqlGetTOTPSecretByUsername: fmt.Sprintf("SELECT secret FROM %s WHERE username=?", totpSecretsTableName),
			sqlUpsertTOTPSecret:        fmt.Sprintf("REPLACE INTO %s (username, secret) VALUES (?, ?)", totpSecretsTableName),

			sqlGetU2FDeviceHandleByUsername: fmt.Sprintf("SELECT keyHandle, publicKey FROM %s WHERE username=?", u2fDeviceHandlesTableName),
			sqlUpsertU2FDeviceHandle:        fmt.Sprintf("REPLACE INTO %s (username, keyHandle, publicKey) VALUES (?, ?, ?)", u2fDeviceHandlesTableName),

			sqlInsertAuthenticationLog:     fmt.Sprintf("INSERT INTO %s (username, successful, time) VALUES (?, ?, ?)", authenticationLogsTableName),
			sqlGetLatestAuthenticationLogs: fmt.Sprintf("SELECT successful, time FROM %s WHERE time>? AND username=? ORDER BY time DESC", authenticationLogsTableName),
		},
	}
	if err := provider.initialize(db); err != nil {
		logging.Logger().Fatalf("Unable to initialize SQLite database %s: %s", path, err)
	}
	return &provider
}
