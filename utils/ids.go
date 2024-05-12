package utils

import "fmt"

func GetID(insertedID interface {}) string  {
  return fmt.Sprintf("%v", insertedID)
}