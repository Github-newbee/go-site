package sid

import (
	"sync"

	"github.com/sony/sonyflake"
)

var (
	sf     *sonyflake.Sonyflake
	sfOnce sync.Once
)

type Sid struct{}

func NewSid() *Sid {
	// 保证只初始化一次,避免NewSid()被多次调用，并发时会出现重复问题
	sfOnce.Do(func() {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{
			MachineID: func() (uint16, error) {
				// 这里可以设置机器ID，确保每个实例的机器ID唯一
				return 1, nil
			},
		})
		if sf == nil {
			panic("sonyflake not created")
		}
	})
	return &Sid{}
}

func (s Sid) GenString() (string, error) {
	id, err := sf.NextID()
	if err != nil {
		return "", err
	}
	return IntToBase62(int(id)), nil
}

func (s Sid) GenUint64() (uint64, error) {
	return sf.NextID()
}
