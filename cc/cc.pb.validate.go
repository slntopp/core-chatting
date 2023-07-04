// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: cc/cc.proto

package cc

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Empty with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Empty) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Empty with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in EmptyMultiError, or nil if none found.
func (m *Empty) ValidateAll() error {
	return m.validate(true)
}

func (m *Empty) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return EmptyMultiError(errors)
	}

	return nil
}

// EmptyMultiError is an error wrapping multiple validation errors returned by
// Empty.ValidateAll() if the designated constraints aren't met.
type EmptyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EmptyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EmptyMultiError) AllErrors() []error { return m }

// EmptyValidationError is the validation error returned by Empty.Validate if
// the designated constraints aren't met.
type EmptyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EmptyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EmptyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EmptyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EmptyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EmptyValidationError) ErrorName() string { return "EmptyValidationError" }

// Error satisfies the builtin error interface
func (e EmptyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEmpty.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EmptyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EmptyValidationError{}

// Validate checks the field values on ChatMeta with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ChatMeta) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChatMeta with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ChatMetaMultiError, or nil
// if none found.
func (m *ChatMeta) ValidateAll() error {
	return m.validate(true)
}

func (m *ChatMeta) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Unread

	if all {
		switch v := interface{}(m.GetLastMessage()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ChatMetaValidationError{
					field:  "LastMessage",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ChatMetaValidationError{
					field:  "LastMessage",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLastMessage()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ChatMetaValidationError{
				field:  "LastMessage",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ChatMetaMultiError(errors)
	}

	return nil
}

// ChatMetaMultiError is an error wrapping multiple validation errors returned
// by ChatMeta.ValidateAll() if the designated constraints aren't met.
type ChatMetaMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChatMetaMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChatMetaMultiError) AllErrors() []error { return m }

// ChatMetaValidationError is the validation error returned by
// ChatMeta.Validate if the designated constraints aren't met.
type ChatMetaValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChatMetaValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChatMetaValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChatMetaValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChatMetaValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChatMetaValidationError) ErrorName() string { return "ChatMetaValidationError" }

// Error satisfies the builtin error interface
func (e ChatMetaValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChatMeta.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChatMetaValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChatMetaValidationError{}

// Validate checks the field values on Chat with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Chat) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Chat with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ChatMultiError, or nil if none found.
func (m *Chat) ValidateAll() error {
	return m.validate(true)
}

func (m *Chat) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Uuid

	// no validation rules for Owner

	// no validation rules for Role

	// no validation rules for Created

	if m.Topic != nil {
		// no validation rules for Topic
	}

	if m.Meta != nil {

		if all {
			switch v := interface{}(m.GetMeta()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ChatValidationError{
						field:  "Meta",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ChatValidationError{
						field:  "Meta",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMeta()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ChatValidationError{
					field:  "Meta",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ChatMultiError(errors)
	}

	return nil
}

// ChatMultiError is an error wrapping multiple validation errors returned by
// Chat.ValidateAll() if the designated constraints aren't met.
type ChatMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChatMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChatMultiError) AllErrors() []error { return m }

// ChatValidationError is the validation error returned by Chat.Validate if the
// designated constraints aren't met.
type ChatValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChatValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChatValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChatValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChatValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChatValidationError) ErrorName() string { return "ChatValidationError" }

// Error satisfies the builtin error interface
func (e ChatValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChat.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChatValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChatValidationError{}

// Validate checks the field values on Chats with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Chats) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Chats with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ChatsMultiError, or nil if none found.
func (m *Chats) ValidateAll() error {
	return m.validate(true)
}

func (m *Chats) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetChats() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ChatsValidationError{
						field:  fmt.Sprintf("Chats[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ChatsValidationError{
						field:  fmt.Sprintf("Chats[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ChatsValidationError{
					field:  fmt.Sprintf("Chats[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ChatsMultiError(errors)
	}

	return nil
}

// ChatsMultiError is an error wrapping multiple validation errors returned by
// Chats.ValidateAll() if the designated constraints aren't met.
type ChatsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChatsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChatsMultiError) AllErrors() []error { return m }

// ChatsValidationError is the validation error returned by Chats.Validate if
// the designated constraints aren't met.
type ChatsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChatsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChatsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChatsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChatsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChatsValidationError) ErrorName() string { return "ChatsValidationError" }

// Error satisfies the builtin error interface
func (e ChatsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChats.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChatsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChatsValidationError{}

// Validate checks the field values on Attachment with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Attachment) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Attachment with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AttachmentMultiError, or
// nil if none found.
func (m *Attachment) ValidateAll() error {
	return m.validate(true)
}

func (m *Attachment) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Type

	// no validation rules for Content

	if len(errors) > 0 {
		return AttachmentMultiError(errors)
	}

	return nil
}

// AttachmentMultiError is an error wrapping multiple validation errors
// returned by Attachment.ValidateAll() if the designated constraints aren't met.
type AttachmentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AttachmentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AttachmentMultiError) AllErrors() []error { return m }

// AttachmentValidationError is the validation error returned by
// Attachment.Validate if the designated constraints aren't met.
type AttachmentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AttachmentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AttachmentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AttachmentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AttachmentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AttachmentValidationError) ErrorName() string { return "AttachmentValidationError" }

// Error satisfies the builtin error interface
func (e AttachmentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAttachment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AttachmentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AttachmentValidationError{}

// Validate checks the field values on Message with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Message) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Message with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in MessageMultiError, or nil if none found.
func (m *Message) ValidateAll() error {
	return m.validate(true)
}

func (m *Message) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Uuid

	// no validation rules for Kind

	// no validation rules for Sender

	// no validation rules for Content

	for idx, item := range m.GetAttachments() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, MessageValidationError{
						field:  fmt.Sprintf("Attachments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, MessageValidationError{
						field:  fmt.Sprintf("Attachments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MessageValidationError{
					field:  fmt.Sprintf("Attachments[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Sent

	// no validation rules for Edited

	// no validation rules for UnderReview

	if m.Chat != nil {
		// no validation rules for Chat
	}

	if len(errors) > 0 {
		return MessageMultiError(errors)
	}

	return nil
}

// MessageMultiError is an error wrapping multiple validation errors returned
// by Message.ValidateAll() if the designated constraints aren't met.
type MessageMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MessageMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MessageMultiError) AllErrors() []error { return m }

// MessageValidationError is the validation error returned by Message.Validate
// if the designated constraints aren't met.
type MessageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MessageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MessageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MessageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MessageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MessageValidationError) ErrorName() string { return "MessageValidationError" }

// Error satisfies the builtin error interface
func (e MessageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMessage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MessageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MessageValidationError{}

// Validate checks the field values on Messages with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Messages) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Messages with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MessagesMultiError, or nil
// if none found.
func (m *Messages) ValidateAll() error {
	return m.validate(true)
}

func (m *Messages) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetMessages() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, MessagesValidationError{
						field:  fmt.Sprintf("Messages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, MessagesValidationError{
						field:  fmt.Sprintf("Messages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MessagesValidationError{
					field:  fmt.Sprintf("Messages[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return MessagesMultiError(errors)
	}

	return nil
}

// MessagesMultiError is an error wrapping multiple validation errors returned
// by Messages.ValidateAll() if the designated constraints aren't met.
type MessagesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MessagesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MessagesMultiError) AllErrors() []error { return m }

// MessagesValidationError is the validation error returned by
// Messages.Validate if the designated constraints aren't met.
type MessagesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MessagesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MessagesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MessagesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MessagesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MessagesValidationError) ErrorName() string { return "MessagesValidationError" }

// Error satisfies the builtin error interface
func (e MessagesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMessages.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MessagesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MessagesValidationError{}

// Validate checks the field values on User with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *User) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on User with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UserMultiError, or nil if none found.
func (m *User) ValidateAll() error {
	return m.validate(true)
}

func (m *User) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Uuid

	// no validation rules for Title

	if len(errors) > 0 {
		return UserMultiError(errors)
	}

	return nil
}

// UserMultiError is an error wrapping multiple validation errors returned by
// User.ValidateAll() if the designated constraints aren't met.
type UserMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserMultiError) AllErrors() []error { return m }

// UserValidationError is the validation error returned by User.Validate if the
// designated constraints aren't met.
type UserValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserValidationError) ErrorName() string { return "UserValidationError" }

// Error satisfies the builtin error interface
func (e UserValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUser.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserValidationError{}

// Validate checks the field values on Defaults with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Defaults) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Defaults with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DefaultsMultiError, or nil
// if none found.
func (m *Defaults) ValidateAll() error {
	return m.validate(true)
}

func (m *Defaults) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DefaultsMultiError(errors)
	}

	return nil
}

// DefaultsMultiError is an error wrapping multiple validation errors returned
// by Defaults.ValidateAll() if the designated constraints aren't met.
type DefaultsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DefaultsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DefaultsMultiError) AllErrors() []error { return m }

// DefaultsValidationError is the validation error returned by
// Defaults.Validate if the designated constraints aren't met.
type DefaultsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DefaultsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DefaultsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DefaultsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DefaultsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DefaultsValidationError) ErrorName() string { return "DefaultsValidationError" }

// Error satisfies the builtin error interface
func (e DefaultsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDefaults.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DefaultsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DefaultsValidationError{}

// Validate checks the field values on Users with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Users) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Users with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UsersMultiError, or nil if none found.
func (m *Users) ValidateAll() error {
	return m.validate(true)
}

func (m *Users) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetUsers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UsersValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UsersValidationError{
						field:  fmt.Sprintf("Users[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UsersValidationError{
					field:  fmt.Sprintf("Users[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UsersMultiError(errors)
	}

	return nil
}

// UsersMultiError is an error wrapping multiple validation errors returned by
// Users.ValidateAll() if the designated constraints aren't met.
type UsersMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UsersMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UsersMultiError) AllErrors() []error { return m }

// UsersValidationError is the validation error returned by Users.Validate if
// the designated constraints aren't met.
type UsersValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UsersValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UsersValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UsersValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UsersValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UsersValidationError) ErrorName() string { return "UsersValidationError" }

// Error satisfies the builtin error interface
func (e UsersValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUsers.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UsersValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UsersValidationError{}

// Validate checks the field values on Event with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Event) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Event with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in EventMultiError, or nil if none found.
func (m *Event) ValidateAll() error {
	return m.validate(true)
}

func (m *Event) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Type

	switch v := m.Item.(type) {
	case *Event_Chat:
		if v == nil {
			err := EventValidationError{
				field:  "Item",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetChat()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, EventValidationError{
						field:  "Chat",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, EventValidationError{
						field:  "Chat",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetChat()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return EventValidationError{
					field:  "Chat",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Event_Msg:
		if v == nil {
			err := EventValidationError{
				field:  "Item",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetMsg()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, EventValidationError{
						field:  "Msg",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, EventValidationError{
						field:  "Msg",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMsg()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return EventValidationError{
					field:  "Msg",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return EventMultiError(errors)
	}

	return nil
}

// EventMultiError is an error wrapping multiple validation errors returned by
// Event.ValidateAll() if the designated constraints aren't met.
type EventMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EventMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EventMultiError) AllErrors() []error { return m }

// EventValidationError is the validation error returned by Event.Validate if
// the designated constraints aren't met.
type EventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EventValidationError) ErrorName() string { return "EventValidationError" }

// Error satisfies the builtin error interface
func (e EventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EventValidationError{}
