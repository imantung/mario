package mario

// partial represents a partial template
type partial struct {
	name   string
	source string
	tpl    *Template
}

// newPartial instanciates a new partial
func newPartial(name string, source string, tpl *Template) *partial {
	return &partial{
		name:   name,
		source: source,
		tpl:    tpl,
	}
}

// template returns parsed partial template
func (p *partial) template() (*Template, error) {
	if p.tpl == nil {
		var err error

		p.tpl, err = New().Parse(p.source)
		if err != nil {
			return nil, err
		}
	}

	return p.tpl, nil
}
