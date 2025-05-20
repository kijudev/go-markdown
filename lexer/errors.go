package lexer

import "errors"

const (
	ErrKindRuneRead   string = "LEXER_ERR_RUNE_READ"
	ErrKindRuneUnread string = "LEXER_ERR_RUNE_UNREAD"
	ErrKindUknown     string = "LEXER_ERR_UKNOWN"
)

var (
	ErrRuneRead   = errors.New(ErrKindRuneRead)
	ErrRuneUnread = errors.New(ErrKindRuneUnread)
	ErrUknown     = errors.New(ErrKindUknown)
)
