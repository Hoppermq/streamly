package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEngineErrorOrder              = errors.New("engine should be set before server")
	ErrEmptyContent                  = errors.New("content could not be empty")
	ErrBatchSizeMaxSizeExceeded      = errors.New("batch size exceeded")
	ErrFailedToReadFile              = errors.New("failed to read file")
	ErrFailedToCompileJSONSchema     = errors.New("failed to compile json schema")
	ErrFailedToReadJSONSchema        = errors.New("failed to read json schema file")
	ErrFailedToUnmarshalJSONSchema   = errors.New("failed to unmarshal json schema")
	ErrFailedToAddJSONSchemaResource = errors.New("failed to add json schema reousrce")
	ErrRootUserIDNotSet              = errors.New("ROOT_USER_ID environment variable not set")
	ErrDatabaseConnection            = errors.New("failed to connect to database")
	ErrDatabaseConstraint            = errors.New("database constraint violation")
)

func FailedToReadJSONSchema(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToReadJSONSchema, err)
}

func FailedToUnmarshalJSONSchema(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToUnmarshalJSONSchema, err)
}

func FailedToCompileJSONSchema(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToCompileJSONSchema, err)
}

func FailedToAddJsonSchemaResource(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToAddJSONSchemaResource, err)
}

func FailedToReadFile(path string) error {
	return fmt.Errorf("%w: %s", ErrFailedToReadFile, path)
}

var (
	ErrTenantIDRequired   = errors.New("tenant_id is required")
	ErrSourceIDRequired   = errors.New("source_id is required")
	ErrMessageIDRequired  = errors.New("message_id is required")
	ErrTopicRequired      = errors.New("topic is required")
	ErrEventTypeRequired  = errors.New("event_type is required")
	ErrRawContentRequired = errors.New("raw_content is required")
	ErrEventEmpty         = errors.New("event cannot be empty")
	ErrEventSize          = errors.New("event size cannot be greater than ~4GB")
)

func EventMessageMissing(eventID int) error {
	return fmt.Errorf("%w: %d", ErrMessageIDRequired, eventID)
}

func EventTypeMissing(eventID int) error {
	return fmt.Errorf("%w: %d", ErrEventTypeRequired, eventID)
}

func EventContentEmpty(eventID int) error {
	return fmt.Errorf("%w: %d", ErrEventEmpty, eventID)
}

var (
	ErrZitadelClientCreation = errors.New("failed to create Zitadel client")
	ErrZitadelPATRequired    = errors.New("pat token is required (use WithPAT or WithPATFromFile)")
)

var (
	ErrSerializerInvalidTimeWindow    = errors.New("invalid time window")
	ErrSerializerInvalidGroupByClause = errors.New("invalid group by clause")
	ErrSerializerInvalidGroupBy       = errors.New("groupBy must be string or time window object")
	ErrSerializerInvalidSelectClause  = errors.New("invalid select clause")
	ErrSerializerInvalidSelect        = errors.New("select must be string or aggregation object")
)

func SerializerInvalidTimeWindow(err error) error {
	return fmt.Errorf("%w: %w", ErrSerializerInvalidTimeWindow, err)
}

func SerializerInvalidSelectFunction(err error) error {
	return fmt.Errorf("%w: %w", ErrSerializerInvalidSelect, err)
}

var (
	ErrFromTranslationFailed    = errors.New("failed to translate FROM")
	ErrSelectTranslationFailed  = errors.New("failed to translate SELECT")
	ErrWhereTranslationFailed   = errors.New("failed to translate WHERE")
	ErrGroupByTranslationFailed = errors.New("failed to translate GROUP BY")
	ErrOrderByTranslationFailed = errors.New("failed to translate ORDER BY")
	ErrSelectClauseEmpty        = errors.New("SELECT clause cannot be empty")
	ErrSelectClauseType         = errors.New("unknown SELECT clause type")
	ErrFromEmpty                = errors.New("FROM datasource cannot be empty")
	ErrInOperatorValue          = errors.New("IN operator requires array value for field: ")
	ErrUnknownGroupBy           = errors.New("unknown GROUP BY clause type")
)

func TranslatorFailedToTranslate(kwError, err error) error {
	return fmt.Errorf("%w: %w", kwError, err)
}

func TranslatorInOperatorInvalidValue(field string) error {
	return fmt.Errorf("%w: %s", ErrInOperatorValue, field)
}

var (
	ErrNoSelectClauseDefined = errors.New("no SELECT clause defined")
	ErrNoFromSourceDefined   = errors.New("no FROM source defined")
	ErrInOperator            = errors.New("IN operator requires []any value")
)

var (
	ErrNilUserInput         = errors.New("user input cannot be nil")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationDelete   = errors.New("failed to delete organization")
	ErrOrganizationUpdate   = errors.New("failed to update organization")
	ErrOrganizationCreate   = errors.New("failed to create organization")
	ErrUserQuery            = errors.New("failed to query user")
	ErrOrganizationQuery    = errors.New("failed to query organization")
)

func OrganizationDeleteFailed(err error) error {
	return fmt.Errorf("%w: %w", ErrOrganizationDelete, err)
}
func OrganizationUpdateFailed(err error) error {
	return fmt.Errorf("%w: %w", ErrOrganizationUpdate, err)
}
func OrganizationCreateFailed(err error) error {
	return fmt.Errorf("%w: %w", ErrOrganizationCreate, err)
}

func RootUserQueryFailed(err error, userID string) error {
	return fmt.Errorf("%w: %w. user: %s", ErrUserQuery, err, userID)
}

func OrganizationQueryFailed(err error, orgID string) error {
	return fmt.Errorf("%w: %w. orgID: %s", ErrOrganizationQuery, err, orgID)
}
