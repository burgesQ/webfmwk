package webfmwk

type (
	ErrorHandled interface {
		GetOPCode() int
		GetContent() interface{}
	}

	Error struct {
		Op      int
		Content interface{}
	}

	Error500 struct {
		content interface{}
	}

	Error400 struct {
		content interface{}
	}

	Error404 struct {
		content interface{}
	}

	Error406 struct {
		content interface{}
	}

	Error422 struct {
		content interface{}
	}

	ErrorUnprocessableEntity Error422
)

func (e Error) GetOPCode() int {
	return e.Op
}

func (e Error) GetContent() interface{} {
	return e.Content
}

func New400(content interface{}) Error400 {
	return Error400{content}
}

func (e Error400) GetOPCode() int {
	return 400
}

func (e Error400) GetContent() interface{} {
	return e.content
}

func New404(content interface{}) Error404 {
	return Error404{content}
}

func (e Error404) GetOPCode() int {
	return 404
}

func (e Error404) GetContent() interface{} {
	return e.content
}

func New406(content interface{}) Error406 {
	return Error406{content}
}

func (e Error406) GetOPCode() int {
	return 406
}

func (e Error406) GetContent() interface{} {
	return e.content
}

func New500(content interface{}) Error500 {
	return Error500{content}
}

func (e Error500) GetOPCode() int {
	return 500
}

func (e Error500) GetContent() interface{} {
	return e.content
}

func New422(content interface{}) Error422 {
	return Error422{content}
}

func (e Error422) GetOPCode() int {
	return 422
}

func (e Error422) GetContent() interface{} {
	return e.content
}
