package lexer

import "errors"

const (
	ErrKindEOF    string = "LEXER_EOF"
	ErrKindRead   string = "LEXER_READ"
	ErrKindUnread string = "LEXER_UNREAD"
)

var (
	ErrEOF    = errors.New(ErrKindEOF)
	ErrRead   = errors.New(ErrKindRead)
	ErrUnread = errors.New(ErrKindUnread)
)
