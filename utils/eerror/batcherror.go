package errorx

import "bytes"

type BatchError []error

func (be BatchError) Error() string {
	var buf bytes.Buffer

	for i := range be {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(be[i].Error())
	}

	return buf.String()
}
