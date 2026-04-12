package utils

import (
	"fmt"
	"time"
)

func GenerateStudentCode() string {
	return fmt.Sprintf("STD-%d", time.Now().UnixNano())
}
