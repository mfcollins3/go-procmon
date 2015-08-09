/*
The MIT License (MIT)

Copyright (c) 2015 Michael F. Collins, III

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including but without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package procmon

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("nullProcessMonitor", func() {
	const testMessage = "Hello, World!"
	var writer = &nullProcessMonitor{}

	Describe("Write", func() {
		var length int
		var err error

		BeforeEach(func() {
			length, err =
				writer.Write(bytes.NewBufferString(testMessage).Bytes())
		})

		It("returns the length of the string", func() {
			Expect(length).To(Equal(len(testMessage)))
		})

		It("does not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("WriteString", func() {
		var length int
		var err error

		BeforeEach(func() {
			length, err = io.WriteString(writer, testMessage)
		})

		It("returns the length of the string", func() {
			Expect(length).To(Equal(len(testMessage)))
		})

		It("does not return an error", func() {
			Expect(err).To(BeNil())
		})
	})
})
