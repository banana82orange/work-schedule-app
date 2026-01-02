package handler

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseTime(t string) *timestamppb.Timestamp {
	if t == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed)
}