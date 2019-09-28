package webfmwk

type (
	ErrorHandled interface {
		GetOPCode() int
		GetContent() interface{}
		IsJSON() bool
	}

	Error struct {
		Op      int
		Content interface{}
		JSON    bool
	}

	Error500 struct {
		content interface{}
	}

	Error404 struct {
		content interface{}
	}
)

func (e Error) GetOPCode() int {
	return e.Op
}

func (e Error) GetContent() interface{} {
	return e.Content
}

func (e Error) IsJSON() bool {
	return e.JSON
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

func (e Error404) IsJSON() bool {
	return true
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

func (e Error500) IsJSON() bool {
	return true
}
