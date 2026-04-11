package util

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
