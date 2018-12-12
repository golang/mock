package vendor_pkg

//go:generate mockgen -destination mock.go -package vendor_pkg test.src/b Ifc
