package gotemplater

type Template struct {
	ID string //each teample has an id, the value is absolute path of that template file
}

func (template *Template) Render(context *Context) (string, error) {
	return "", nil
}
