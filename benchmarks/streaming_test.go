package pixy

import (
	"bytes"
	"testing"
)

func renderIcon() string {
	b := acquireBytesBuffer()
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	return b.String()
}

func streamIcon(b *bytes.Buffer) {
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
}

func render() string {
	b := acquireBytesBuffer()
	b.WriteString(renderIcon())
	b.WriteString(renderIcon())
	b.WriteString(renderIcon())
	b.WriteString(renderIcon())
	return b.String()
}

func stream() string {
	b := acquireBytesBuffer()
	streamIcon(b)
	streamIcon(b)
	streamIcon(b)
	streamIcon(b)
	return b.String()
}

func BenchmarkHTMLRendering(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			render()
		}
	})
}

func BenchmarkHTMLStreaming(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stream()
		}
	})
}
