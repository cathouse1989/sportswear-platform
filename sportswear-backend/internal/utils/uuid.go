package utils

import "github.com/google/uuid"

// StringToUUID 将字符串转换为 UUID
func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// StringPtrToUUIDPtr 将 *string 转换为 *uuid.UUID
func StringPtrToUUIDPtr(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	uid, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &uid
}
