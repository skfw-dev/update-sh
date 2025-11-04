package common

import (
	"encoding/base64"
	"fmt"

	"go.skfw.net/shikacore/caseconv"
)

type Field struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func NewField(name string, value any) *Field {
	return &Field{
		Key:   name,
		Value: value,
	}
}

func (s *Field) GetKey() string {
	return caseconv.ToCamelCase(s.Key)
}

func (s *Field) GetValue() any {
	return s.Value
}

func (s *Field) ToString() string {
	switch v := s.Value.(type) {
	case string:
		return v
	case []byte:
		return base64.StdEncoding.EncodeToString(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (s *Field) ToMap() map[string]any {
	return map[string]any{
		s.Key: s.Value,
	}
}

type Fields []Field

func (s Fields) Count() int {
	return len(s)
}

func (s Fields) ToMap() map[string]any {
	result := make(map[string]any)
	for _, field := range s {
		result[(&field).GetKey()] = field.Value
	}
	return result
}
