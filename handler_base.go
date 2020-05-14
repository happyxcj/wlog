package wlog


type BaseHandler struct {
	w       Writer
	encoder Encoder
}

func NewBaseHandler(w Writer, encoder Encoder) *BaseHandler {
	return &BaseHandler{
		w:       w,
		encoder: encoder,
	}
}

func (h *BaseHandler) With(fields ...Field) Handler {
	return NewWithHandler(h).With(fields...)
}

func (h *BaseHandler) Write(entry *Entry, fields ...Field) error {
	buf := GetBuf()
	err := h.encoder.Encode(buf, entry, fields...)
	if err != nil {
		PutBuf(buf)
		return err
	}
	_, err = h.w.Write(buf.Bytes())
	PutBuf(buf)
	return err
}

func (h *BaseHandler) Flush() error {
	return h.w.Flush()
}

func (h *BaseHandler) Close() error {
	return h.w.Close()
}