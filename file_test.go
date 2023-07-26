package kcli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getCurrDir", func() {
	var (
		dir string

		path string
	)
	JustBeforeEach(func() {
		dir = getCurrDir(path)
	})
	Context("pathが空文字の時", func() {
		BeforeEach(func() {
			path = ""
		})
		It("dirも空文字にになる", func() {
			Expect(dir).To(Equal(""))
		})
	})
	Context("ちゃんとしたpathが渡された時", func() {
		BeforeEach(func() {
			path = "/dir1/dir2/dir3"
		})
		It("カレントディクレトリはdir3になる", func() {
			Expect(dir).To(Equal("dir3"))
		})
	})
})
