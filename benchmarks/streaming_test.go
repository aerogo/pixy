package pixy

import (
	"bytes"
	"strings"
	"testing"
)

func renderIconBuilder() string {
	b := acquireStringsBuilder()
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	return b.String()
}

func renderIcon() string {
	b := acquireStringsBuilder()
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	return b.String()
}

func streamIconBuilder(b *strings.Builder) {
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
}

func streamIcon(b *bytes.Buffer) {
	b.WriteString("<icon name='test'></icon>")
	b.WriteString("<icon name='test'></icon>")
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

func renderBuilder() string {
	b := acquireStringsBuilder()
	b.WriteString(renderIconBuilder())
	b.WriteString(renderIconBuilder())
	b.WriteString(renderIconBuilder())
	b.WriteString(renderIconBuilder())
	return b.String()
}

func streamBuilder() string {
	b := acquireStringsBuilder()
	streamIconBuilder(b)
	streamIconBuilder(b)
	streamIconBuilder(b)
	streamIconBuilder(b)
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

func BenchmarkHTMLRenderingBuilder(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			renderBuilder()
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

func BenchmarkHTMLStreamingBuilder(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			streamBuilder()
		}
	})
}
