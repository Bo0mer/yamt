package flagvar

import (
	"fmt"
	"strings"
)

type MapFlag map[string]string

func (m *MapFlag) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MapFlag) Set(value string) error {
	if *m == nil {
		*m = make(map[string]string)
	}

	kv := strings.Split(value, "=")
	if len(kv) != 2 || len(kv[0]) == 0 || len(kv[1]) == 0 {
		return fmt.Errorf("unsupported map flag format: %q", value)
	}
	key, value := kv[0], kv[1]
	(*m)[key] = value
	return nil
}
