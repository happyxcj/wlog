package wlog

type WithHandler struct {
	Handler
	fields []Field
}

func NewWithHandler(inner Handler) *WithHandler {
	return &WithHandler{
		Handler:      inner,
		fields: make([]Field, 0, 2),
	}
}

func (h *WithHandler) With(fields ...Field) Handler {
	h.fields = append(h.fields, fields...)
	return h
}

func (h *WithHandler) Write(entry *Entry, fields ...Field) error {
	if len(h.fields) > 0 {
		fields = append(h.fields, fields...)
	}
	return h.Handler.Write(entry, fields...)
}
