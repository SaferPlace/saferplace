package sqldatabase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"api.safer.place/incident/v1"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"safer.place/realtime/internal/database"
)

// Database contains the database connection
type Database struct {
	db *sql.DB

	hasIncidentStmt            *sql.Stmt
	saveIncidentStmt           *sql.Stmt
	updateResolutionStmt       *sql.Stmt
	saveCommentStmt            *sql.Stmt
	viewIncidentStmt           *sql.Stmt
	viewCommentsStmt           *sql.Stmt
	incidentsWithoutReviewStmt *sql.Stmt
}

// New creates a new SQL database
func New() (*Database, error) {
	var cfg Config
	if err := envconfig.Process("SQL", &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse SQL config: %w", err)
	}
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// TODO: Probably uncomment to allow migrations
	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("unable to prepare database: %w", err)
	}

	hasIncidentStmt, err := db.Prepare("SELECT id FROM incidents WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare hasIncidents query: %w", err)
	}
	saveIncidentStmt, err := db.Prepare(saveIncidentQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare saveIncident query: %w", err)
	}
	updateResolutionStmt, err := db.Prepare(updateResolutionQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare updateResolution query: %w", err)
	}
	saveCommentStmt, err := db.Prepare(saveCommentQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare saveComment query: %w", err)
	}
	viewIncidentStmt, err := db.Prepare(viewIncidentQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare viewIncident query: %w", err)
	}
	viewCommentsStmt, err := db.Prepare(viewCommentsQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare viewComments query: %w", err)
	}
	incidentsWithoutReviewStmt, err := db.Prepare(incidentsWithoutReviewQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare incidentsWithoutReview query: %w", err)
	}

	return &Database{
		db:                         db,
		hasIncidentStmt:            hasIncidentStmt,
		saveIncidentStmt:           saveIncidentStmt,
		updateResolutionStmt:       updateResolutionStmt,
		saveCommentStmt:            saveCommentStmt,
		viewIncidentStmt:           viewIncidentStmt,
		viewCommentsStmt:           viewCommentsStmt,
		incidentsWithoutReviewStmt: incidentsWithoutReviewStmt,
	}, nil
}

// SaveIncident to the sql database
func (db *Database) SaveIncident(ctx context.Context, inc *incident.Incident) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if exists, err := db.hasIncident(ctx, tx, inc.Id); err != nil {
		return fmt.Errorf("unable to check does the incident exist: %w", err)
	} else if exists {
		return database.ErrAlreadyExists
	}

	if _, err := tx.Stmt(db.saveIncidentStmt).ExecContext(ctx,
		inc.Id,
		inc.Timestamp,
		inc.Description,
		inc.Lat,
		inc.Lon,
		inc.Resolution.String(),
	); err != nil {
		return fmt.Errorf("unable to save incident: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

// SaveReview updates the incident record with the resolution and adds a comment.
func (db *Database) SaveReview(
	ctx context.Context,
	id string,
	res incident.Resolution,
	comment *incident.Comment,
) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if exists, err := db.hasIncident(ctx, tx, id); err != nil {
		return fmt.Errorf("unable to check does the incident exist: %w", err)
	} else if !exists {
		return database.ErrDoesNotExist
	}

	if _, err := tx.Stmt(db.updateResolutionStmt).ExecContext(
		ctx, res.String(), id,
	); err != nil {
		return fmt.Errorf("unable to update incident resolution: %w", err)
	}

	if _, err := tx.Stmt(db.saveCommentStmt).ExecContext(
		ctx,
		uuid.New().String(), // id
		id,                  // incident_id
		comment.Timestamp,   // timestamp
		comment.AuthorId,    // author
		comment.Message,     // comment
		res.String(),        // resolution

	); err != nil {
		return fmt.Errorf("unable to save comment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

// ViewIncident recovers incident information
func (db *Database) ViewIncident(ctx context.Context, id string) (*incident.Incident, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	inc := new(incident.Incident)
	resStr := ""
	if err := tx.Stmt(db.viewIncidentStmt).QueryRow(id).Scan(
		&inc.Id,
		&inc.Timestamp,
		&inc.Description,
		&inc.Lat,
		&inc.Lon,
		&resStr,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrDoesNotExist
		}
		return nil, fmt.Errorf("unable to get incident info: %w", err)
	}
	inc.Resolution = incident.Resolution(incident.Resolution_value[resStr])

	// TODO: Get comments
	rows, err := tx.Stmt(db.viewCommentsStmt).QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get incident comments: %w", err)
	}
	for rows.Next() {
		comment := new(incident.Comment)
		discard := ""
		if err := rows.Scan(
			&discard,
			&discard,
			&comment.Timestamp, // timestamp
			&comment.AuthorId,  // author
			&comment.Message,   // comment
			&discard,           // resolution, TODO: Maybe add to comment?
		); err != nil {
			return nil, fmt.Errorf("unable to scan comment: %w", err)
		}
		inc.ReviewerComments = append(inc.ReviewerComments, comment)
	}

	// Sort the reviewer comments
	sort.Sort(ByTimestamp(inc.ReviewerComments))

	// We don't actually change anything but we are using this to close the
	// transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %w", err)
	}
	return inc, nil
}

// IncidentsWithoutReview gets all the incidents which have the UNDEFINED
func (db *Database) IncidentsWithoutReview(ctx context.Context) ([]*incident.Incident, error) {
	rows, err := db.incidentsWithoutReviewStmt.QueryContext(ctx,
		incident.Resolution_RESOLUTION_UNSPECIFIED.String(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*incident.Incident{}, nil
		}
		return nil, fmt.Errorf("unable list incidents: %w", err)
	}

	incidents := make([]*incident.Incident, 0)
	for rows.Next() {
		inc := new(incident.Incident)
		resStr := ""
		if err := rows.Scan(
			&inc.Id,
			&inc.Timestamp,
			&inc.Description,
			&inc.Lat,
			&inc.Lon,
			&resStr,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, database.ErrDoesNotExist
			}
			return nil, fmt.Errorf("unable to get incident info: %w", err)
		}
		inc.Resolution = incident.Resolution(incident.Resolution_value[resStr])
		incidents = append(incidents, inc)
	}

	return incidents, nil
}

func (db *Database) hasIncident(ctx context.Context, tx *sql.Tx, id string) (bool, error) {
	// First check do we already have an entry, so we can return already exists
	row := tx.Stmt(db.hasIncidentStmt).QueryRowContext(ctx, id)

	var existingID string
	if err := row.Scan(&existingID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("unable to check for existance of incident: %w", err)
	}

	return id == existingID, nil
}

// ByTimestamp sorts the comments by timestamp, from the oldest to the newest
type ByTimestamp []*incident.Comment

func (s ByTimestamp) Len() int           { return len(s) }
func (s ByTimestamp) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByTimestamp) Less(i, j int) bool { return s[i].Timestamp < s[j].Timestamp }

// Config of the SQLDatabase
type Config struct {
	Driver string `default:"sqlite3"`
	DSN    string `default:"file:incidents.db"`
}

var createTableQuery = `
CREATE TABLE IF NOT EXISTS incidents (
	id TEXT PRIMARY KEY,
	timestamp INTEGER NOT NULL,
	description TEXT,
	lat REAL NOT NULL,
	lon REAL NOT NULL,
	resolution TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
	id TEXT PRIMARY KEY,
	incident_id TEXT NOT NULL,
	timestamp INTEGER NOT NULL,
	author TEXT NOT NULL,
	comment TEXT NOT NULL,
	resolution TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS incident_ids ON comments (incident_id);
`

var saveIncidentQuery = `
INSERT INTO incidents
	(id, timestamp, description, lat, lon, resolution)
VALUES
	(?, ?, ?, ?, ?, ?);
`

var updateResolutionQuery = `
UPDATE incidents
SET
	resolution=?
WHERE
	id=?;
`

var saveCommentQuery = `
INSERT INTO comments
	(id, incident_id, timestamp, author, comment, resolution)
VALUES
	(?, ?, ?, ?, ?, ?);
`

var viewIncidentQuery = `
SELECT * FROM incidents WHERE id=?;
`

var viewCommentsQuery = `
SELECT * FROM comments WHERE incident_id=?;
`

var incidentsWithoutReviewQuery = `
SELECT * FROM incidents WHERE resolution=?;
`
