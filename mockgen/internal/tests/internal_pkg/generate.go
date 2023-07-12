package test

//go:generate mockgen -destination subdir/internal/pkg/reflect_output/mock.go go.uber.org/mock/mockgen/internal/tests/internal_pkg/subdir/internal/pkg Intf
//go:generate mockgen -source subdir/internal/pkg/input.go -destination subdir/internal/pkg/source_output/mock.go
