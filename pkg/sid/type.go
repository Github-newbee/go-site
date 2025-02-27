package sid

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SnowflakeID int64

// 在将结构体序列化为 JSON 时自动调用
func (s SnowflakeID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, s)), nil
}

// 在将 JSON 反序列化为结构体时自动调用。
func (s *SnowflakeID) UnmarshalJSON(data []byte) error {
	idStr := strings.Trim(string(data), `"`)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}
	*s = SnowflakeID(id)
	return nil
}

// 在从数据库查询结果中读取值时自动调用。
func (s *SnowflakeID) Scan(value interface{}) error {
	if v, ok := value.(int64); ok {
		*s = SnowflakeID(v)
		return nil
	}
	return errors.New("failed to scan SnowflakeID")
}

// 在将值存储到数据库中时自动调用。
func (s SnowflakeID) Value() (driver.Value, error) {
	return int64(s), nil
}

// NewSnowflakeIDFromString converts a string to SnowflakeID
// 将类型转换成当前类型
func NewSnowflakeIDFromString(idStr string) (SnowflakeID, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return SnowflakeID(id), nil
}
