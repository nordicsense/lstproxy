module github.com/nordicsense/lstproxy

go 1.15

// Uncomment for MODIS and GDAL development
// replace github.com/nordicsense/modis => ../modis
// replace github.com/nordicsense/gdal => ../gdal

require (
	github.com/gonum/blas v0.0.0-20181208220705-f22b278b28ac // indirect; indirect (matching)
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82 // indirect
	github.com/gonum/internal v0.0.0-20181124074243-f884aa714029 // indirect
	github.com/gonum/lapack v0.0.0-20181123203213-e4cdc5a0bff9 // indirect
	github.com/gonum/matrix v0.0.0-20181209220409-c518dec07be9
	github.com/nordicsense/modis v0.0.0-20210621221057-a2ec700d5c4f
	github.com/schollz/progressbar/v2 v2.15.0
	github.com/subchen/go-xmldom v1.1.2 // indirect (finmeteo)
	github.com/teris-io/cli v1.0.1
)
