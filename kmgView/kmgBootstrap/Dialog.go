package kmgBootstrap
import (
	"github.com/bronze1man/kmg/kmgView"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgXss"
)

type DialogOpenWithButton struct{
	ButtonTitle string
	DialogContent kmgView.HtmlRenderer
}

func (config DialogOpenWithButton) HtmlRender()string{
	targetId:="dialog_"+kmgRand.MustCryptoRandToReadableAlphaNum(10)
	return `<button type="button" class="btn btn-primary btn-lg" data-toggle="modal" data-target="#`+targetId+`">
`+kmgXss.H(config.ButtonTitle)+`
</button>
<div class="modal fade" id="`+targetId+`" tabindex="-1" role="dialog">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
      </div>
      <div class="modal-body">
        `+config.DialogContent.HtmlRender()+`
      </div>
    </div>
  </div>
</div>
`
}