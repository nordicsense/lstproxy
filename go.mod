module github.com/nordicsense/lstproxy

go 1.15

replace github.com/lukeroth/gdal => github.com/osklyar/gdal v0.0.0-20210506223057-2fecd07bf2d8

require (
	github.com/gonum/blas v0.0.0-20181208220705-f22b278b28ac // indirect (matching)
	github.com/lukeroth/gdal v0.0.0-20210429104814-0ed96ec28fb2
	github.com/nordicsense/modis v0.0.0-20210509214422-ca11fe4bf0d6
	github.com/schollz/progressbar/v2 v2.15.0
	github.com/subchen/go-xmldom v1.1.2 // indirect (finmeteo)
	github.com/teris-io/cli v1.0.1
)
