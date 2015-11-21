package kmgBootstrap

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

func tplForm(config Form) string {
	method := "post"
	if config.IsGet {
		method = "get"
	}
	var _buf bytes.Buffer
	_buf.WriteString(`    <form `)
	if config.IsHidden {
		_buf.WriteString(`style="display: none;"`)
	}
	_buf.WriteString(` id="`)
	_buf.WriteString(kmgXss.H(config.Id))
	_buf.WriteString(`" class="form-horizontal" autocomplete="off" role="form" action="`)
	_buf.WriteString(kmgXss.H(config.Url))
	_buf.WriteString(`" method="`)
	_buf.WriteString(kmgXss.H(method))
	_buf.WriteString(`">
        <div class="panel-body">
            `)
	for _, input := range config.InputList {
		_buf.WriteString(`                `)
		_buf.WriteString(input.HtmlRender())
		_buf.WriteString(`            `)
	}
	_buf.WriteString(`        </div>

    `)
	if !config.NoSubmit {
		_buf.WriteString(`        <div class="panel-footer">
            <center>
                <button type="submit" class="btn btn-primary" style="width: 170px;">
                    <span class="fa fa-check"></span>
                    提交
                </button>
            </center>
        </div>
    `)
	}
	_buf.WriteString(`    </form>`)
	return _buf.String()
}
