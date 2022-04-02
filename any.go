//go:build !go1.18
// +build !go1.18

package i18n

// 添加「编译约束」，参考： https://pkg.go.dev/cmd/go#hdr-Build_constraints
type any = interface{}
