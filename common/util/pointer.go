package util

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// transforms string literal to string pointer, without this function an inline string pointer cannot be created
func InlineStringPointer(s string) *string {
	return &s
}

// transforms unsigned 16 byte int literal to unsigned 16 byte int pointer, without this function an inline unsigned 16 byte int pointer cannot be created
func InlineUInt16Pointer(i uint16) *uint16 {
	return &i
}

// transforms unsigned 32 byte int literal to unsigned 32 byte int pointer, without this function an inline unsigned 32 byte int pointer cannot be created
func InlineUInt32Pointer(i uint32) *uint32 {
	return &i
}

func ProtoStringValue(s *string) *wrapperspb.StringValue {
	if s != nil {
		return wrapperspb.String(*s)
	}
	return nil
}

func ProtoUInt32Value(ui *uint32) *wrapperspb.UInt32Value {
	if ui != nil {
		return wrapperspb.UInt32(*ui)
	}
	return nil
}
