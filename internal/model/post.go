package model

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/gofrs/uuid"
)

// PostEntry ...
type PostEntry struct {
	UUID       uuid.UUID `db:"uuid"`
	AuthorUUID uuid.UUID `db:"author_uuid"`
	Posts      []byte    `db:"posts"` // CBOR serialized
	CreatedAt  uint64    // unix ts
}

// Posts are not stored directly in the db,
// they are always serialized inside a post entry via CBOR
type Post struct {
	PostEntryUUID uuid.UUID // UUID of the originating post entry for this content
	AuthorUUID    uuid.UUID
	CreatedAt     uint64 // unix ts
	Content       string // json encoded rich content? markdown?
}

// NewEntry ...
func NewEntry(posts []Post, authorUUID uuid.UUID) (*PostEntry, error) {
	postUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	serializedPosts, err := cbor.Marshal(posts)
	if err != nil {
		return nil, err
	}

	return &PostEntry{
		UUID:       postUUID,
		AuthorUUID: authorUUID,
		Posts:      serializedPosts,
	}, nil
}

// GetPosts ...
func (p *PostEntry) GetPosts() ([]Post, error) {
	var res []Post
	err := cbor.Unmarshal(p.Posts, res)
	return res, err
}
