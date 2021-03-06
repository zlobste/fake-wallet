// +build !skippackr
// Code generated by github.com/gobuffalo/packr/v2. DO NOT EDIT.

// You can use the "packr2 clean" command to clean up this,
// and any other packr generated files.
package packrd

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file/resolver"
)

var _ = func() error {
	const gk = "f571d1ba9f65b89d1bd69d62ef12d9f7"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"51692660bf4af56300ac7511807c78ed": "1f8b08000000000000ff3cca410ac2301005d07d4ef17769309522b872a7281ea01e20951f0d368964a682b717147cebd7f758e5746b4189cbd3a4226c8a54b4228850059dbcf354678f12323d22e9cc2bcc0b059d8dc7f16c3dec293c08ea9d8d4bb61ec37ad83a6ff0d3d9b81f0fff3725bdd654be6de376e6130000fffff158d9a684000000",
		"aa4639612eafee59c35a8463df12c078": "1f8b08000000000000ff9454cdce9b3010bcf3147b04f58b4453f51435a7be42cfd5626f122bfe41bb26347dfa0a6208101af1edc9d2ccb03bc3dabb1d7c71e6cc18097ed59962ea4e112b4bd008b16479060060343caa3267213668a166e390ef70a53bf810c137d67ef4648f8e7af20d595d90f3ef6501cf9a93a5e1077f0b991c1abbf5cb358ab481f548de9753f640ce8a4336378e221417ceb7f896bbab821ddb7d7d335b8a688b8d13f561fac6111b957ffb807db1c1448bd63e5da0d64c220091fec4998565cd9bf751fc36baf36f7c04a613317945925282dce8e25557a145af6832f6be00752175857cc08e3fa02c16bad07ae2f57efd3a8eedde388f8c5e504513fcc6f5fdbf7f21af8921e5369926c50b794a7669844991b9117f5a882e343ec26a72093b42b91279da935561870daa176134e9ba760789e86a684dbc84263ea0bfc1d3daba4ddf8e9fa1f599e6504fdf0e50280a351da6485a9b356888660d9bfed491f02f0000ffff3fbc5c58be040000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("migrations", "./migrations")
		b.SetResolver("001-initial.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "aa4639612eafee59c35a8463df12c078"})
		b.SetResolver("002-assets.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "51692660bf4af56300ac7511807c78ed"})
	}()
	return nil
}()
