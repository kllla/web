package ascii

import (
	"fmt"
	"github.com/kllla/web/src/common/sao"
	"github.com/mbndr/figlet4go"
)

var ascii, options = initASCII()

func initASCII() (*figlet4go.AsciiRender, *figlet4go.RenderOptions) {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colours to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontName = "bloody"
	bucket := sao.NewSao()
	bin := bucket.GetStaticFiles("Bloody.flf")
	ascii.LoadBindataFont(bin,"bloody")
	return ascii, options
}

// RenderString returns the string in the format defined in initASCII
func RenderString(text string) string {
	fmt.Printf("Rendering Banner %s\n", text)
	renderStr, _ := ascii.RenderOpts(text, options)
	return renderStr
}

