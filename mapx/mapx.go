package mapx

import (
	"encoding/xml"
	"io"
)

type Map map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	//构建xml输出头部
	var err error
	for key, value := range m {
		name := xml.Name{Space: "", Local: key}
		err = e.EncodeToken(xml.StartElement{Name: name})
		if err != nil {
			return err
		}
		err = e.EncodeToken(xml.CharData(value))
		if err != nil {
			return err
		}
		err = e.EncodeToken(xml.EndElement{Name: name})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Map{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

/**
 * @desc 校验
 */
func IsExist(data map[string]interface{}, name string) bool {
	_, ok := data[name]
	return ok
}
