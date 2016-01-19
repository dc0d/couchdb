package couchdb

import (
	"fmt"
	"strings"
)

// Server gives access to the welcome string and version information.
// http://docs.couchdb.org/en/latest/intro/api.html#server
type Server struct {
	Couchdb string
	UUID    string
	Vendor  struct {
		Version string
		Name    string
	}
	Version string
}

// DatabaseInfo has info about the specified database.
// http://docs.couchdb.org/en/latest/api/database/common.html#get--db
type DatabaseInfo struct {
	DbName             string `json:"db_name"`
	DocCount           int    `json:"doc_count"`
	DocDelCount        int    `json:"doc_del_count"`
	UpdateSeq          int    `json:"update_seq"`
	PurgeSeq           int    `json:"purge_seq"`
	CompactRunning     bool   `json:"compact_running"`
	DiskSize           int    `json:"disk_size"`
	DataSize           int    `json:"data_size"`
	InstanceStartTime  string `json:"instance_start_time"`
	DiskFormatVersion  int    `json:"disk_format_version"`
	CommittedUpdateSeq int    `json:"committed_update_seq"`
}

// DatabaseResponse is body for successful database calls.
type DatabaseResponse struct {
	Ok bool
}

// Error describes CouchDB error.
type Error struct {
	Method     string
	URL        string
	StatusCode int
	Type       string `json:"error"`
	Reason     string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"CouchDB - %s %s, Status Code: %d, Error: %s, Reason: %s",
		e.Method,
		e.URL,
		e.StatusCode,
		e.Type,
		e.Reason,
	)
}

// CouchDoc describes interface for every couchdb document.
// GetDocument() is just a dummy method for now to satisfy the interface.
type CouchDoc interface {
	GetDocument() *Document
}

// Document is base struct which should be embedded by any other couchdb document.
type Document struct {
	ID          string                `json:"_id,omitempty"`
	Rev         string                `json:"_rev,omitempty"`
	Attachments map[string]Attachment `json:"_attachments,omitempty"`
}

// GetDocument is just a dummy method to describe the CouchDoc interface
func (d *Document) GetDocument() *Document {
	return d
}

// DesignDocument is a special type of CouchDB document that contains application code.
// http://docs.couchdb.org/en/latest/json-structure.html#design-document
type DesignDocument struct {
	Document
	Language string                        `json:"language,omitempty"`
	Views    map[string]DesignDocumentView `json:"views,omitempty"`
}

// Name returns design document name without the "_design/" prefix
func (dd DesignDocument) Name() string {
	return strings.TrimPrefix(dd.ID, "_design/")
}

// DesignDocumentView contains map/reduce functions.
type DesignDocumentView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

// Attachment uploads a document in multipart/related format.
// http://docs.couchdb.org/en/latest/api/document/common.html#creating-multiple-attachments
type Attachment struct {
	Follows     bool   `json:"follows"`
	ContentType string `json:"content_type"`
	Length      int64  `json:"length"`
}

// DocumentResponse is response for multipart/related file upload.
type DocumentResponse struct {
	Ok  bool
	ID  string
	Rev string
}

// Task describes currently running task.
// http://docs.couchdb.org/en/latest/api/server/common.html#active-tasks
type Task struct {
	ChangesDone  int `json:"changes_done"`
	Database     string
	Pid          string
	Progress     int
	StartedOn    int `json:"started_on"`
	Status       string
	Task         string
	TotalChanges int `json:"total_changes"`
	Type         string
	UpdatedOn    string `json:"updated_on"`
}

// QueryParameters is struct to define url query parameters for design documents.
// http://docs.couchdb.org/en/latest/api/ddoc/views.html#db-design-design-doc-view-view-name
type QueryParameters struct {
	Conflicts       *bool   `url:"conflicts,omitempty"`
	Descending      *bool   `url:"descending,omitempty"`
	Group           *bool   `url:"group,omitempty"`
	IncludeDocs     *bool   `url:"include_docs,omitempty"`
	Attachments     *bool   `url:"attachments,omitempty"`
	AttEncodingInfo *bool   `url:"att_encoding_info,omitempty"`
	InclusiveEnd    *bool   `url:"inclusive_end,omitempty"`
	Reduce          *bool   `url:"reduce,omitempty"`
	UpdateSeq       *bool   `url:"update_seq,omitempty"`
	GroupLevel      *int32  `url:"group_level,omitempty"`
	Limit           *int32  `url:"limit,omitempty"`
	Skip            *int32  `url:"skip,omitempty"`
	Key             *string `url:"key,omitempty"`
	EndKey          *string `url:"endkey,comma,omitempty"`
	EndKeyDocID     *string `url:"end_key_doc_id,omitempty"`
	Stale           *string `url:"stale,omitempty"`
	StartKey        *string `url:"startkey,comma,omitempty"`
	StartKeyDocID   *string `url:"startkey_docid,omitempty"`
}

// NewQueryParameters returns query parameters with default values
// http://docs.couchdb.org/en/latest/api/ddoc/views.html#get--db-_design-ddoc-_view-view
// The problem is "reduce" for example. The default value is true.
// If we have a map/reduce function that has a reduce part everything works as expected.
// We'll get into trouble if we want to reuse this document without reduce.
// If we use the omitempty flag on the Reduce field it would get it's zero value false
// which would not be sent to the server.
// // func NewQueryParameters() QueryParameters {
// // 	// reduce is the exception. the default would be true
// // 	// but as have have more cases where we don't have a reduce function we set it false
// // 	// set it to true if you really need it.
// // 	return QueryParameters{
// // 		Conflicts:       false,
// // 		Descending:      false,
// // 		Group:           false,
// // 		IncludeDocs:     false,
// // 		Attachments:     false,
// // 		AttEncodingInfo: false,
// // 		InclusiveEnd:    true,
// // 		Reduce:          false,
// // 		Skip:            0,
// // 		UpdateSeq:       false,
// // 	}
// // }

// ViewResponse is response for querying design documents.
type ViewResponse struct {
	Offset    int   `json:"offset,omitempty"`
	Rows      []Row `json:"rows,omitempty"`
	TotalRows int   `json:"total_rows,omitempty"`
	UpdateSeq int   `json:"update_seq,omitempty"`
}

// Row is single row inside design document query response.
type Row struct {
	ID    string                 `json:"id"`
	Key   interface{}            `json:"key"`
	Value interface{}            `json:"value,omitempty"`
	Doc   map[string]interface{} `json:"doc,omitempty"`
}

// BulkDoc describes POST /db/_bulk_docs request object.
// http://docs.couchdb.org/en/latest/api/database/bulk-api.html#post--db-_bulk_docs
type BulkDoc struct {
	AllOrNothing bool          `json:"all_or_nothing,omitempty"`
	NewEdits     bool          `json:"new_edits,omitempty"`
	Docs         []interface{} `json:"docs"`
}

// Credentials has information about POST _session form parameters.
// http://docs.couchdb.org/en/latest/api/server/authn.html#cookie-authentication
type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// PostSessionResponse is response from posting to session api.
type PostSessionResponse struct {
	Ok    bool
	Name  string
	Roles []string
}

// User is special CouchDB document format.
// http://docs.couchdb.org/en/latest/intro/security.html#users-documents
type User struct {
	Document
	DerivedKey     string   `json:"derived_key,omitempty"`
	Name           string   `json:"name,omitempty"`
	Roles          []string `json:"roles"`
	Password       string   `json:"password,omitempty"`     // plain text password when creating the user
	PasswordSha    string   `json:"password_sha,omitempty"` // hashed password when requesting user information
	PasswordScheme string   `json:"password_scheme,omitempty"`
	Salt           string   `json:"salt,omitempty"`
	Type           string   `json:"type,omitempty"`
	Iterations     int      `json:"iterations,omitempty"`
}

// NewUser returns new user instance.
func NewUser(name, password string, roles []string) User {
	user := User{
		Document: Document{
			ID: "org.couchdb.user:" + name,
		},
		DerivedKey:     "",
		Name:           name,
		Roles:          roles,
		Password:       password,
		PasswordSha:    "",
		PasswordScheme: "",
		Salt:           "",
		Type:           "user",
	}
	return user
}

// GetSessionResponse returns complete information about authenticated user.
// http://docs.couchdb.org/en/latest/api/server/authn.html#get--_session
type GetSessionResponse struct {
	Info struct {
		Authenticated          string   `json:"authenticated"`
		AuthenticationDb       string   `json:"authentication_db"`
		AuthenticationHandlers []string `json:"authentication_handlers"`
	} `json:"info"`
	Ok          bool `json:"ok"`
	UserContext struct {
		Db    string   `json:"db"`
		Name  string   `json:"name"`
		Roles []string `json:"roles"`
	} `json:"userCtx"`
}
