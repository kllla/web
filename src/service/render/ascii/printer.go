package ascii

import (
	"github.com/kllla/web/src/common/sao"
	"github.com/mbndr/figlet4go"
)

var ascii, options = initASCII()

func initASCII() (*figlet4go.AsciiRender, *figlet4go.RenderOptions) {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colours to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontName = "ansiregular"
	bucket := sao.New()
	bin := bucket.GetStaticFiles("ansiregular.flf")
	ascii.LoadBindataFont(bin, "ansiregular")
	return ascii, options
}

// RenderString returns the string in the format defined in initASCII
func RenderString(text string) string {
	renderStr, _ := ascii.RenderOpts(text, options)
	return renderStr
}
