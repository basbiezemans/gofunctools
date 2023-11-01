module github.com/basbiezemans/gofunctools

go 1.21

retract (
	v1.5.2 // Backwards incompatible with v1.5.1
	v1.5.1 // Backwards incompatible with v1.5
)

require golang.org/x/exp v0.0.0-20231006140011-7918f672742d
