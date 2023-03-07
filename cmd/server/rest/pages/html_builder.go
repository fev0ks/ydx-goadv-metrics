package pages

import "bytes"

type htmlBuilder struct {
	bytes.Buffer
}

// GetHTMLBuilder - конструктор html страницы
func GetHTMLBuilder() *htmlBuilder {
	builder := htmlBuilder{}
	builder.WriteString(OHtml)
	return &builder
}

func (hb *htmlBuilder) GetHTMLPage() string {
	hb.WriteString(CHtml)
	return hb.String()
}

func (hb *htmlBuilder) AddHeader(text string) *htmlBuilder {
	hb.WriteString(OH2)
	hb.WriteString(text)
	hb.WriteString(CH2)
	return hb
}

func (hb *htmlBuilder) Add(text string) *htmlBuilder {
	hb.WriteString(text)
	return hb
}
